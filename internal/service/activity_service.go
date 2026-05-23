package service

import (
	"time"

	"github.com/Hayversong/questboard/internal/model"
)

func AddActivity(
	project *model.Project,
	message string,
) {

	project.Activities =
		append(
			[]model.Activity{
				{
					Message: message,
					Time: time.Now().Format(
						"02/01 15:04",
					),
				},
			},
			project.Activities...,
		)

	if len(project.Activities) > 20 {

		project.Activities =
			project.Activities[:20]

	}
}
