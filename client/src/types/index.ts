export interface User {
  id: string;
  username: string;
  email: string;
}

export interface Project {
  id: string;
  name: string;
  owner: User;
  tasks: Task[];
}

export interface Task {
  id: string;
  title: string;
  status: string;
  createdAt: string;
  project: Project;
}

export interface AuthResponse {
  token: string;
  user: User;
}
