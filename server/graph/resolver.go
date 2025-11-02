package graph

import (
	"task-manager/apps/projects"
	"task-manager/apps/tasks"
	"task-manager/apps/users"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	userSrv *users.UserService
	projectSrv *projects.ProjectService
	taskSrv *tasks.TaskService
}

func NewResolver(userSrv *users.UserService, 
	projectSrv *projects.ProjectService, 
	taskSrv *tasks.TaskService) *Resolver {
		return &Resolver{
			userSrv: userSrv,
			projectSrv: projectSrv,
			taskSrv: taskSrv,
		}
	}
