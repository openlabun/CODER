package roble_infrastructure

import (
	"context"
	"fmt"
	"strings"
	"time"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
)

const (
	challengeTableName = "Challenge"
)

type ChallengeRepository struct {
	adapter *infrastructure.RobleDatabaseAdapter
}

func NewChallengeRepository(adapter *infrastructure.RobleDatabaseAdapter) *ChallengeRepository {
	return &ChallengeRepository{adapter: adapter}
}

func (r *ChallengeRepository) CreateChallenge(ctx context.Context, challenge *Entities.Challenge) (*Entities.Challenge, error) {
	if challenge == nil {
		return nil, fmt.Errorf("challenge is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	if err := r.upsertChallengeIOVariables(ctx, challenge); err != nil {
		return nil, err
	}

	_, err := r.adapter.Insert(challengeTableName, []map[string]any{challengeToRecord(challenge)})
	if err != nil {
		return nil, err
	}

	return r.GetChallengeByID(ctx, challenge.ID)
}

func (r *ChallengeRepository) UpdateChallenge(ctx context.Context, challenge *Entities.Challenge) (*Entities.Challenge, error) {
	if challenge == nil {
		return nil, fmt.Errorf("challenge is nil")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	challengeID := strings.TrimSpace(challenge.ID)
	if challengeID == "" {
		return nil, fmt.Errorf("challenge id is required")
	}

	if err := r.upsertChallengeIOVariables(ctx, challenge); err != nil {
		return nil, err
	}

	updates := challengeToUpdates(challenge)
	updates["UpdatedAt"] = time.Now().UTC().Format(time.RFC3339)

	_, err := r.adapter.Update(challengeTableName, "ID", challengeID, updates)
	if err != nil {
		return nil, err
	}

	return r.GetChallengeByID(ctx, challengeID)
}

func (r *ChallengeRepository) DeleteChallenge(ctx context.Context, challengeID string) error {
	normalizedID := strings.TrimSpace(challengeID)
	if normalizedID == "" {
		return fmt.Errorf("challengeID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return err
	}

	res, err := r.adapter.Read(challengeTableName, map[string]string{"ID": normalizedID})
	if err != nil {
		return err
	}

	var ioVariableIDs []string
	if record, findErr := firstRecord(res); findErr == nil {
		ioVariableIDs = relatedIOVariableIDs(record)
	}

	_, err = r.adapter.Delete(challengeTableName, "ID", normalizedID)
	if err != nil {
		return err
	}

	return deleteIOVariablesByIDs(ctx, r.adapter, ioVariableIDs)
}

func (r *ChallengeRepository) GetChallengeByID(ctx context.Context, challengeID string) (*Entities.Challenge, error) {
	normalizedID := strings.TrimSpace(challengeID)
	if normalizedID == "" {
		return nil, fmt.Errorf("challengeID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(challengeTableName, map[string]string{"ID": normalizedID})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return r.recordToHydratedChallenge(ctx, record)
}

func (r *ChallengeRepository) GetChallenges(ctx context.Context, status, tag, difficulty *string) ([]*Entities.Challenge, error) {
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	filters := map[string]string{}
	if status != nil {
		if s := strings.TrimSpace(*status); s != "" {
			filters["Status"] = s
		}
	}
	if tag != nil {
		if t := strings.TrimSpace(*tag); t != "" {
			filters["Tags"] = t
		}
	}
	if difficulty != nil {
		if d := strings.TrimSpace(*difficulty); d != "" {
			filters["Difficulty"] = d
		}
	}

	var res map[string]any
	var err error
	if len(filters) == 0 {
		res, err = r.adapter.Read(challengeTableName, nil)
	} else {
		res, err = r.adapter.Read(challengeTableName, filters)
	}
	if err != nil {
		return nil, err
	}

	records := extractRecords(res)
	if len(records) == 0 {
		return []*Entities.Challenge{}, nil
	}

	challenges := make([]*Entities.Challenge, 0, len(records))
	for _, record := range records {
		challenge, mapErr := r.recordToHydratedChallenge(ctx, record)
		if mapErr != nil {
			return nil, mapErr
		}
		if challenge != nil {
			challenges = append(challenges, challenge)
		}
	}

	return challenges, nil
}

func (r *ChallengeRepository) GetChallengesByUserID(ctx context.Context, userID string, examID *string) ([]*Entities.Challenge, error) {
	normalizedUserID := strings.TrimSpace(userID)

	if normalizedUserID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	// 1. Si examID es proporcionado, obtener los challenges relacionados a ese examID y userID
	if examID != nil {
		examChallenges, err := r.GetChallengesByExamID(ctx, *examID)
		if err != nil {
			return nil, err
		}

		return examChallenges, nil
	}

	// 2. Si examID no es proporcionado, obtener todos los challenges relacionados al userID
	res, err := r.adapter.Read(challengeTableName, map[string]string{"UserID": normalizedUserID})
	if err != nil {
		return nil, err
	}

	records := extractRecords(res)
	challenges := make([]*Entities.Challenge, 0, len(records))
	for _, record := range records {
		challenge, mapErr := r.recordToHydratedChallenge(ctx, record)
		if mapErr != nil {
			return nil, mapErr
		}
		if challenge != nil {
			challenges = append(challenges, challenge)
		}
	}

	return challenges, nil
}

func (r *ChallengeRepository) GetChallengesByExamID(ctx context.Context, examID string) ([]*Entities.Challenge, error) {
	normalizedID := strings.TrimSpace(examID)
	if normalizedID == "" {
		return nil, fmt.Errorf("examID is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	// 1. Obtener ExamItems por ExamID
	examItemRes, err := r.adapter.Read("ExamItem", map[string]string{"ExamID": normalizedID})
	if err != nil {
		return nil, err
	}
	examItemRecords := extractRecords(examItemRes)
	if len(examItemRecords) == 0 {
		return []*Entities.Challenge{}, nil
	}

	// 2. Extraer ChallengeIDs únicos
	challengeIDSet := make(map[string]struct{})
	for _, item := range examItemRecords {
		id := asString(item["ChallengeID"])
		if id != "" {
			challengeIDSet[id] = struct{}{}
		}
	}
	if len(challengeIDSet) == 0 {
		return []*Entities.Challenge{}, nil
	}

	// 3. Consultar Challenge por cada ChallengeID
	challenges := make([]*Entities.Challenge, 0, len(challengeIDSet))
	for challengeID := range challengeIDSet {
		res, err := r.adapter.Read(challengeTableName, map[string]string{"ID": challengeID})
		if err != nil {
			return nil, err
		}
		records := extractRecords(res)
		for _, record := range records {
			challenge, mapErr := r.recordToHydratedChallenge(ctx, record)
			if mapErr != nil {
				return nil, mapErr
			}
			if challenge != nil {
				challenges = append(challenges, challenge)
			}
		}
	}
	return challenges, nil
}

func (r *ChallengeRepository) GetChallengesByTag(ctx context.Context, tag string) ([]*Entities.Challenge, error) {
	normalizedTag := strings.TrimSpace(tag)
	if normalizedTag == "" {
		return nil, fmt.Errorf("tag is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(challengeTableName, map[string]string{"Tags": normalizedTag})
	if err != nil {
		return nil, err
	}

	records := extractRecords(res)
	if len(records) == 0 {
		return []*Entities.Challenge{}, nil
	}

	challenges := make([]*Entities.Challenge, 0, len(records))
	for _, record := range records {
		challenge, mapErr := r.recordToHydratedChallenge(ctx, record)
		if mapErr != nil {
			return nil, mapErr
		}
		if challenge != nil {
			challenges = append(challenges, challenge)
		}
	}

	return challenges, nil
}

func (r *ChallengeRepository) GetInputVariablesByChallengeID(ctx context.Context, challengeID string) ([]*Entities.IOVariable, error) {
	challenge, err := r.GetChallengeByID(ctx, challengeID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return []*Entities.IOVariable{}, nil
	}

	variables := make([]*Entities.IOVariable, 0, len(challenge.InputVariables))
	for i := range challenge.InputVariables {
		variable := challenge.InputVariables[i]
		variables = append(variables, &variable)
	}

	return variables, nil
}

func (r *ChallengeRepository) GetOutputVariablesByChallengeID(ctx context.Context, challengeID string) ([]*Entities.IOVariable, error) {
	challenge, err := r.GetChallengeByID(ctx, challengeID)
	if err != nil {
		return nil, err
	}
	if challenge == nil || strings.TrimSpace(challenge.OutputVariable.ID) == "" {
		return []*Entities.IOVariable{}, nil
	}

	output := challenge.OutputVariable
	return []*Entities.IOVariable{&output}, nil
}

func (r *ChallengeRepository) recordToHydratedChallenge(ctx context.Context, record map[string]any) (*Entities.Challenge, error) {
	inputIDs := asStringList(record["Input"])
	if len(inputIDs) == 0 {
		inputIDs = asStringList(record["InputVariables"])
	}

	outputID := asString(record["Output"])
	if strings.TrimSpace(outputID) == "" {
		outputID = asString(record["OutputVariable"])
	}

	inputVariables, err := r.getIOVariablesByIDs(ctx, inputIDs)
	if err != nil {
		return nil, err
	}

	outputVariable, err := r.getIOVariableByID(ctx, outputID)
	if err != nil {
		return nil, err
	}

	return recordToChallenge(record, inputVariables, outputVariable)
}

func (r *ChallengeRepository) upsertChallengeIOVariables(ctx context.Context, challenge *Entities.Challenge) error {
	for _, input := range challenge.InputVariables {
		if err := r.upsertIOVariable(ctx, input); err != nil {
			return err
		}
	}

	if err := r.upsertIOVariable(ctx, challenge.OutputVariable); err != nil {
		return err
	}

	return nil
}

func (r *ChallengeRepository) getIOVariablesByIDs(ctx context.Context, ids []string) ([]Entities.IOVariable, error) {
	if len(ids) == 0 {
		return []Entities.IOVariable{}, nil
	}

	variables := make([]Entities.IOVariable, 0, len(ids))
	for _, id := range ids {
		ioVariable, err := r.getIOVariableByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if ioVariable != nil {
			variables = append(variables, *ioVariable)
		}
	}

	return variables, nil
}

func (r *ChallengeRepository) getIOVariableByID(ctx context.Context, variableID string) (*Entities.IOVariable, error) {
	normalizedID := strings.TrimSpace(variableID)
	if normalizedID == "" {
		return nil, nil
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return nil, err
	}

	res, err := r.adapter.Read(ioVariableTableName, map[string]string{"ID": normalizedID})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return recordToIOVariable(record)
}

func (r *ChallengeRepository) upsertIOVariable(ctx context.Context, variable Entities.IOVariable) error {
	variableID := strings.TrimSpace(variable.ID)
	if variableID == "" {
		return fmt.Errorf("io variable id is required")
	}
	if err := infrastructure.SetAdapterTokenFromContext(ctx, r.adapter); err != nil {
		return err
	}

	res, err := r.adapter.Read(ioVariableTableName, map[string]string{"ID": variableID})
	if err != nil {
		return err
	}

	if _, findErr := firstRecord(res); findErr != nil {
		_, insertErr := r.adapter.Insert(ioVariableTableName, []map[string]any{ioVariableToRecord(variable)})
		return insertErr
	}

	_, updateErr := r.adapter.Update(ioVariableTableName, "ID", variableID, ioVariableToUpdates(variable))
	return updateErr
}
