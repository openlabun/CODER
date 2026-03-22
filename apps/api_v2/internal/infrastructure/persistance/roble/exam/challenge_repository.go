package roble_infrastructure

import (
	"fmt"
	"strings"
	"time"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
)

const (
	challengeTableName  = "Challenge"
)

type ChallengeRepository struct {
	adapter *infrastructure.RobleDatabaseAdapter
}

func NewChallengeRepository(adapter *infrastructure.RobleDatabaseAdapter) *ChallengeRepository {
	return &ChallengeRepository{adapter: adapter}
}

func (r *ChallengeRepository) CreateChallenge(challenge *Entities.Challenge) (*Entities.Challenge, error) {
	if challenge == nil {
		return nil, fmt.Errorf("challenge is nil")
	}

	if err := r.upsertChallengeIOVariables(challenge); err != nil {
		return nil, err
	}

	_, err := r.adapter.Insert(challengeTableName, []map[string]any{challengeToRecord(challenge)})
	if err != nil {
		return nil, err
	}

	return r.GetChallengeByID(challenge.ID)
}

func (r *ChallengeRepository) UpdateChallenge(challenge *Entities.Challenge) (*Entities.Challenge, error) {
	if challenge == nil {
		return nil, fmt.Errorf("challenge is nil")
	}

	challengeID := strings.TrimSpace(challenge.ID)
	if challengeID == "" {
		return nil, fmt.Errorf("challenge id is required")
	}

	if err := r.upsertChallengeIOVariables(challenge); err != nil {
		return nil, err
	}

	updates := challengeToUpdates(challenge)
	updates["UpdatedAt"] = time.Now().UTC().Format(time.RFC3339)

	_, err := r.adapter.Update(challengeTableName, "ID", challengeID, updates)
	if err != nil {
		return nil, err
	}

	return r.GetChallengeByID(challengeID)
}

func (r *ChallengeRepository) DeleteChallenge(challengeID string) error {
	normalizedID := strings.TrimSpace(challengeID)
	if normalizedID == "" {
		return fmt.Errorf("challengeID is required")
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

	return deleteIOVariablesByIDs(r.adapter, ioVariableIDs)
}

func (r *ChallengeRepository) GetChallengeByID(challengeID string) (*Entities.Challenge, error) {
	normalizedID := strings.TrimSpace(challengeID)
	if normalizedID == "" {
		return nil, fmt.Errorf("challengeID is required")
	}

	res, err := r.adapter.Read(challengeTableName, map[string]string{"ID": normalizedID})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return r.recordToHydratedChallenge(record)
}

func (r *ChallengeRepository) GetChallengesByExamID(examID string) ([]*Entities.Challenge, error) {
	normalizedID := strings.TrimSpace(examID)
	if normalizedID == "" {
		return nil, fmt.Errorf("examID is required")
	}

	res, err := r.adapter.Read(challengeTableName, map[string]string{"ExamID": normalizedID})
	if err != nil {
		return nil, err
	}

	records := extractRecords(res)
	if len(records) == 0 {
		return []*Entities.Challenge{}, nil
	}

	challenges := make([]*Entities.Challenge, 0, len(records))
	for _, record := range records {
		challenge, mapErr := r.recordToHydratedChallenge(record)
		if mapErr != nil {
			return nil, mapErr
		}
		if challenge != nil {
			challenges = append(challenges, challenge)
		}
	}

	return challenges, nil
}

func (r *ChallengeRepository) GetChallengesByTag(tag string) ([]*Entities.Challenge, error) {
	normalizedTag := strings.TrimSpace(tag)
	if normalizedTag == "" {
		return nil, fmt.Errorf("tag is required")
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
		challenge, mapErr := r.recordToHydratedChallenge(record)
		if mapErr != nil {
			return nil, mapErr
		}
		if challenge != nil {
			challenges = append(challenges, challenge)
		}
	}

	return challenges, nil
}

func (r *ChallengeRepository) GetInputVariablesByChallengeID(challengeID string) ([]*Entities.IOVariable, error) {
	challenge, err := r.GetChallengeByID(challengeID)
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

func (r *ChallengeRepository) GetOutputVariablesByChallengeID(challengeID string) ([]*Entities.IOVariable, error) {
	challenge, err := r.GetChallengeByID(challengeID)
	if err != nil {
		return nil, err
	}
	if challenge == nil || strings.TrimSpace(challenge.OutputVariable.ID) == "" {
		return []*Entities.IOVariable{}, nil
	}

	output := challenge.OutputVariable
	return []*Entities.IOVariable{&output}, nil
}

func (r *ChallengeRepository) recordToHydratedChallenge(record map[string]any) (*Entities.Challenge, error) {
	inputIDs := asStringList(record["Input"])
	if len(inputIDs) == 0 {
		inputIDs = asStringList(record["InputVariables"])
	}

	outputID := asString(record["Output"])
	if strings.TrimSpace(outputID) == "" {
		outputID = asString(record["OutputVariable"])
	}

	inputVariables, err := r.getIOVariablesByIDs(inputIDs)
	if err != nil {
		return nil, err
	}

	outputVariable, err := r.getIOVariableByID(outputID)
	if err != nil {
		return nil, err
	}

	return recordToChallenge(record, inputVariables, outputVariable)
}

func (r *ChallengeRepository) upsertChallengeIOVariables(challenge *Entities.Challenge) error {
	for _, input := range challenge.InputVariables {
		if err := r.upsertIOVariable(input); err != nil {
			return err
		}
	}

	if err := r.upsertIOVariable(challenge.OutputVariable); err != nil {
		return err
	}

	return nil
}

func (r *ChallengeRepository) getIOVariablesByIDs(ids []string) ([]Entities.IOVariable, error) {
	if len(ids) == 0 {
		return []Entities.IOVariable{}, nil
	}

	variables := make([]Entities.IOVariable, 0, len(ids))
	for _, id := range ids {
		ioVariable, err := r.getIOVariableByID(id)
		if err != nil {
			return nil, err
		}
		if ioVariable != nil {
			variables = append(variables, *ioVariable)
		}
	}

	return variables, nil
}

func (r *ChallengeRepository) getIOVariableByID(variableID string) (*Entities.IOVariable, error) {
	normalizedID := strings.TrimSpace(variableID)
	if normalizedID == "" {
		return nil, nil
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

func (r *ChallengeRepository) upsertIOVariable(variable Entities.IOVariable) error {
	variableID := strings.TrimSpace(variable.ID)
	if variableID == "" {
		return fmt.Errorf("io variable id is required")
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
