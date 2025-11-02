package tasks

import (
    "context"
    "fmt"
    "strconv"
    "time"
    
    "task-manager/ent"
    "task-manager/ent/project"
    "task-manager/ent/task"
    "task-manager/ent/user"
    "task-manager/graph/model"
    "task-manager/shared/middleware"
)

type TaskService struct {
    client *ent.Client
}

func NewTaskService(client *ent.Client) *TaskService {
    return &TaskService{client: client}
}

func (s *TaskService) CreateTask(ctx context.Context, title string, projectID string) (*model.Task, error) {

    userID, ok := middleware.GetUserID(ctx)
    if !ok {
        return nil, fmt.Errorf("usuario no autenticado")
    }

    projID, err := strconv.Atoi(projectID)
    if err != nil {
        return nil, fmt.Errorf("ID de proyecto inválido: %w", err)
    }
    
    projectExists, err := s.client.Project.
        Query().
        Where(
            project.IDEQ(projID),
            project.HasOwnerWith(user.IDEQ(userID)),
        ).
        Exist(ctx)
    
    if err != nil {
        return nil, fmt.Errorf("error al verificar proyecto: %w", err)
    }
    
    if !projectExists {
        return nil, fmt.Errorf("proyecto no encontrado o no tienes permiso")
    }
    
    entTask, err := s.client.Task.
        Create().
        SetTitle(title).
        SetProjectID(projID).
        SetStatus("pending").
        SetCreatedAt(time.Now()).
        Save(ctx)
    
    if err != nil {
        return nil, fmt.Errorf("error al crear tarea: %w", err)
    }
    
    proj, err := entTask.QueryProject().WithOwner().Only(ctx)
    if err != nil {
        return nil, fmt.Errorf("error al obtener proyecto: %w", err)
    }
    
    owner, _ := proj.Edges.OwnerOrErr()
    
    return &model.Task{
        ID:        strconv.Itoa(entTask.ID),
        Title:     entTask.Title,
        Status:    entTask.Status,
        CreatedAt: entTask.CreatedAt,
        Project: &model.Project{
            ID:   strconv.Itoa(proj.ID),
            Name: proj.Name,
            Owner: &model.User{
                ID:       strconv.Itoa(owner.ID),
                Username: owner.Username,
                Email:    owner.Email,
            },
        },
    }, nil
}


// GetProjectTasks obtiene todas las tareas de un proyecto específico
func (s *TaskService) GetProjectTasks(ctx context.Context, projectID string) ([]*model.Task, error) {
    // Obtener el usuario autenticado
    userID, ok := middleware.GetUserID(ctx)
    if !ok {
        return nil, fmt.Errorf("usuario no autenticado")
    }
    
    // Convertir projectID
    projID, err := strconv.Atoi(projectID)
    if err != nil {
        return nil, fmt.Errorf("ID de proyecto inválido: %w", err)
    }
    
    // Verificar que el proyecto pertenece al usuario
    projectExists, err := s.client.Project.
        Query().
        Where(
            project.IDEQ(projID),
            project.HasOwnerWith(user.IDEQ(userID)),
        ).
        Exist(ctx)
    
    if err != nil {
        return nil, fmt.Errorf("error al verificar proyecto: %w", err)
    }
    
    if !projectExists {
        return nil, fmt.Errorf("proyecto no encontrado o no tienes permiso")
    }
    
    // Obtener las tareas del proyecto
    entTasks, err := s.client.Task.
        Query().
        Where(task.HasProjectWith(project.IDEQ(projID))).
        WithProject(func(q *ent.ProjectQuery) {
            q.WithOwner()
        }).
        All(ctx)
    
    if err != nil {
        return nil, fmt.Errorf("error al obtener tareas: %w", err)
    }
    
    // Convertir a modelo GraphQL
    tasks := make([]*model.Task, len(entTasks))
    for i, t := range entTasks {
        proj, _ := t.Edges.ProjectOrErr()
        owner, _ := proj.Edges.OwnerOrErr()
        
        tasks[i] = &model.Task{
            ID:        strconv.Itoa(t.ID),
            Title:     t.Title,
            Status:    t.Status,
            CreatedAt: t.CreatedAt,
            Project: &model.Project{
                ID:   strconv.Itoa(proj.ID),
                Name: proj.Name,
                Owner: &model.User{
                    ID:       strconv.Itoa(owner.ID),
                    Username: owner.Username,
                    Email:    owner.Email,
                },
            },
        }
    }
    
    return tasks, nil
}

// UpdateTask actualiza una tarea existente
func (s *TaskService) UpdateTask(ctx context.Context, taskID string, title *string, status *string) (*model.Task, error) {
    userID, ok := middleware.GetUserID(ctx)
    if !ok {
        return nil, fmt.Errorf("usuario no autenticado")
    }
    
    id, err := strconv.Atoi(taskID)
    if err != nil {
        return nil, fmt.Errorf("ID de tarea inválido: %w", err)
    }
    
    entTask, err := s.client.Task.
        Query().
        Where(
            task.IDEQ(id),
            task.HasProjectWith(
                project.HasOwnerWith(user.IDEQ(userID)),
            ),
        ).
        Only(ctx)
    
    if err != nil {
        if ent.IsNotFound(err) {
            return nil, fmt.Errorf("tarea no encontrada o no tienes permiso")
        }
        return nil, fmt.Errorf("error al buscar tarea: %w", err)
    }
    
    updateBuilder := entTask.Update()
    
    if title != nil {
        updateBuilder = updateBuilder.SetTitle(*title)
    }
    
    if status != nil {
        updateBuilder = updateBuilder.SetStatus(*status)
    }
    
    _, err = updateBuilder.Save(ctx)
    if err != nil {
        return nil, fmt.Errorf("error al actualizar tarea: %w", err)
    }
    
    entTask, err = s.client.Task.
        Query().
        Where(task.IDEQ(id)).
        WithProject(func(q *ent.ProjectQuery) {
            q.WithOwner()
        }).
        Only(ctx)
    
    if err != nil {
        return nil, fmt.Errorf("error al recargar tarea: %w", err)
    }
    
    proj, _ := entTask.Edges.ProjectOrErr()
    owner, _ := proj.Edges.OwnerOrErr()
    
    return &model.Task{
        ID:        strconv.Itoa(entTask.ID),
        Title:     entTask.Title,
        Status:    entTask.Status,
        CreatedAt: entTask.CreatedAt,
        Project: &model.Project{
            ID:   strconv.Itoa(proj.ID),
            Name: proj.Name,
            Owner: &model.User{
                ID:       strconv.Itoa(owner.ID),
                Username: owner.Username,
                Email:    owner.Email,
            },
        },
    }, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID string) (bool, error) {
    userID, ok := middleware.GetUserID(ctx)
    if !ok {
        return false, fmt.Errorf("usuario no autenticado")
    }
    
    id, err := strconv.Atoi(taskID)
    if err != nil {
        return false, fmt.Errorf("ID de tarea inválido: %w", err)
    }
    
    exists, err := s.client.Task.
        Query().
        Where(
            task.IDEQ(id),
            task.HasProjectWith(
                project.HasOwnerWith(user.IDEQ(userID)),
            ),
        ).
        Exist(ctx)
    
    if err != nil {
        return false, fmt.Errorf("error al verificar tarea: %w", err)
    }
    
    if !exists {
        return false, fmt.Errorf("tarea no encontrada o no tienes permiso")
    }

    err = s.client.Task.DeleteOneID(id).Exec(ctx)
    if err != nil {
        return false, fmt.Errorf("error al eliminar tarea: %w", err)
    }
    
    return true, nil
}