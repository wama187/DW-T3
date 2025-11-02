package projects

import (
	"context"
	"fmt"
	"strconv"

	"task-manager/ent"
	"task-manager/ent/project"
	"task-manager/ent/user"
	"task-manager/graph/model"
	"task-manager/shared/middleware"
)

type ProjectService struct {
    client *ent.Client
}

func NewProjectService(client *ent.Client) *ProjectService {
    return &ProjectService{client: client}
}

func (s *ProjectService) CreateProject(ctx context.Context, name string) (*model.Project, error) {
    userID, ok := middleware.GetUserID(ctx)
    if !ok {
        return nil, fmt.Errorf("usuario no autenticado")
    }
    
    // Crear el proyecto
    entProject, err := s.client.Project.
        Create().
        SetName(name).
        SetOwnerID(userID).
        Save(ctx)
    
    if err != nil {
        return nil, fmt.Errorf("error al crear proyecto: %w", err)
    }
    
    // Obtener el owner para la respuesta
    owner, err := entProject.QueryOwner().Only(ctx)
    if err != nil {
        return nil, fmt.Errorf("error al obtener owner: %w", err)
    }
    
    return &model.Project{
        ID:   strconv.Itoa(entProject.ID),
        Name: entProject.Name,
        Owner: &model.User{
            ID:       strconv.Itoa(owner.ID),
            Username: owner.Username,
            Email:    owner.Email,
        },
        Tasks: []*model.Task{},
    }, nil
}

func (s *ProjectService) GetAllProjects(ctx context.Context) ([]*model.Project, error) {

    userID, ok := middleware.GetUserID(ctx)
    if !ok {
        return nil, fmt.Errorf("usuario no autenticado")
    }
    
    entProjects, err := s.client.Project.
        Query().
        Where(project.HasOwnerWith(user.IDEQ(userID))).
        WithOwner().
        WithTasks().
        All(ctx)
    
    if err != nil {
        return nil, fmt.Errorf("error al obtener proyectos: %w", err)
    }
    
    projects := make([]*model.Project, len(entProjects))
    for i, p := range entProjects {
        owner, err := p.Edges.OwnerOrErr()
        if err != nil {
            return nil, fmt.Errorf("error al obtener owner: %w", err)
        }
        

        var tasks []*model.Task
        if p.Edges.Tasks != nil {
            tasks = make([]*model.Task, len(p.Edges.Tasks))
            for j, t := range p.Edges.Tasks {
                tasks[j] = &model.Task{
                    ID:        strconv.Itoa(t.ID),
                    Title:     t.Title,
                    Status:    t.Status,
                    CreatedAt: t.CreatedAt,
                }
            }
        } else {
            tasks = []*model.Task{}
        }
        
        projects[i] = &model.Project{
            ID:   strconv.Itoa(p.ID),
            Name: p.Name,
            Owner: &model.User{
                ID:       strconv.Itoa(owner.ID),
                Username: owner.Username,
                Email:    owner.Email,
            },
            Tasks: tasks,
        }
    }
    
    return projects, nil
}

func (s *ProjectService) GetProjectByID(ctx context.Context, projectID string) (*model.Project, error) {
    userID, ok := middleware.GetUserID(ctx)
    if !ok {
        return nil, fmt.Errorf("usuario no autenticado")
    }
    
    id, err := strconv.Atoi(projectID)
    if err != nil {
        return nil, fmt.Errorf("ID de proyecto inválido: %w", err)
    }
    
    entProject, err := s.client.Project.
        Query().
        Where(
            project.IDEQ(id),
            project.HasOwnerWith(user.IDEQ(userID)),
        ).
        WithOwner().
        WithTasks().
        Only(ctx)
    
    if err != nil {
        if ent.IsNotFound(err) {
            return nil, fmt.Errorf("proyecto no encontrado o no tienes permiso")
        }
        return nil, fmt.Errorf("error al obtener proyecto: %w", err)
    }
    
    owner, err := entProject.Edges.OwnerOrErr()
    if err != nil {
        return nil, fmt.Errorf("error al obtener owner: %w", err)
    }
    
    var tasks []*model.Task
    if entProject.Edges.Tasks != nil {
        tasks = make([]*model.Task, len(entProject.Edges.Tasks))
        for i, t := range entProject.Edges.Tasks {
            tasks[i] = &model.Task{
                ID:        strconv.Itoa(t.ID),
                Title:     t.Title,
                Status:    t.Status,
                CreatedAt: t.CreatedAt,
            }
        }
    } else {
        tasks = []*model.Task{}
    }
    
    return &model.Project{
        ID:   strconv.Itoa(entProject.ID),
        Name: entProject.Name,
        Owner: &model.User{
            ID:       strconv.Itoa(owner.ID),
            Username: owner.Username,
            Email:    owner.Email,
        },
        Tasks: tasks,
    }, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, projectID string, name string) (*model.Project, error) {
    userID, ok := middleware.GetUserID(ctx)
    if !ok {
        return nil, fmt.Errorf("usuario no autenticado")
    }
    
    id, err := strconv.Atoi(projectID)
    if err != nil {
        return nil, fmt.Errorf("ID de proyecto inválido: %w", err)
    }
    
    entProject, err := s.client.Project.
        Query().
        Where(
            project.IDEQ(id),
            project.HasOwnerWith(user.IDEQ(userID)),
        ).
        Only(ctx)
    
    if err != nil {
        if ent.IsNotFound(err) {
            return nil, fmt.Errorf("proyecto no encontrado o no tienes permiso")
        }
        return nil, fmt.Errorf("error al buscar proyecto: %w", err)
    }
    
    _, err = entProject.Update().
        SetName(name).
        Save(ctx)
    
    if err != nil {
        return nil, fmt.Errorf("error al actualizar proyecto: %w", err)
    }

    entProject, err = s.client.Project.
        Query().
        Where(project.IDEQ(id)).
        WithOwner().
        WithTasks().
        Only(ctx)
    
    if err != nil {
        return nil, fmt.Errorf("error al recargar proyecto: %w", err)
    }
    
    owner, _ := entProject.Edges.OwnerOrErr()
    
    var tasks []*model.Task
    if entProject.Edges.Tasks != nil {
        tasks = make([]*model.Task, len(entProject.Edges.Tasks))
        for i, t := range entProject.Edges.Tasks {
            tasks[i] = &model.Task{
                ID:        strconv.Itoa(t.ID),
                Title:     t.Title,
                Status:    t.Status,
                CreatedAt: t.CreatedAt,
            }
        }
    } else {
        tasks = []*model.Task{}
    }
    
    return &model.Project{
        ID:   strconv.Itoa(entProject.ID),
        Name: entProject.Name,
        Owner: &model.User{
            ID:       strconv.Itoa(owner.ID),
            Username: owner.Username,
            Email:    owner.Email,
        },
        Tasks: tasks,
    }, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, projectID string) (bool, error) {
    userID, ok := middleware.GetUserID(ctx)
    if !ok {
        return false, fmt.Errorf("usuario no autenticado")
    }
    
    id, err := strconv.Atoi(projectID)
    if err != nil {
        return false, fmt.Errorf("ID de proyecto inválido: %w", err)
    }
    
    exists, err := s.client.Project.
        Query().
        Where(
            project.IDEQ(id),
            project.HasOwnerWith(user.IDEQ(userID)),
        ).
        Exist(ctx)
    
    if err != nil {
        return false, fmt.Errorf("error al verificar proyecto: %w", err)
    }
    
    if !exists {
        return false, fmt.Errorf("proyecto no encontrado o no tienes permiso")
    }
    
    err = s.client.Project.DeleteOneID(id).Exec(ctx)
    if err != nil {
        if cErr, ok := err.(*ent.ConstraintError); ok {
            fmt.Print("Error:", cErr)
            return false, fmt.Errorf("no se puede eliminar el proyecto porque tiene tareas asignadas")
        }
        return false, fmt.Errorf("error al eliminar proyecto: %w", err)
    }
    
    return true, nil
}