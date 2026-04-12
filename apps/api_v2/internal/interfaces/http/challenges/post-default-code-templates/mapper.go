package postdefaultcodetemplates

import examEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"

func MapCodeTemplatesToResponse(templates []examEntities.CodeTemplate) map[string]string {
	if len(templates) == 0 {
		return map[string]string{}
	}

	out := make(map[string]string, len(templates))
	for _, tpl := range templates {
		out[string(tpl.Language)] = tpl.Template
	}

	return out
}
