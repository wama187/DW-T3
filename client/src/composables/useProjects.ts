import { useQuery, useMutation } from '@vue/apollo-composable';
import gql from 'graphql-tag';

export const useProjectsQuery = () => {
  return useQuery(gql`
    query projects {
      projects {
        id
        name
      }
    }
  `);
};

export const useProjectQuery = (id: string) => {
  return useQuery(gql`
    query project($id: ID!) {
      project(id: $id) {
        id
        name
        tasks {
          id
          title
          status
        }
      }
    }
  `, { id });
};

export const useCreateProjectMutation = () => {
  return useMutation(gql`
    mutation createProject($name: String!) {
      createProject(name: $name) {
        id
        name
      }
    }
  `);
};

export const useUpdateProjectMutation = () => {
  return useMutation(gql`
    mutation updateProject($id: ID!, $name: String!) {
      updateProject(id: $id, name: $name) {
        id
        name
      }
    }
  `);
};

export const useDeleteProjectMutation = () => {
  return useMutation(gql`
    mutation deleteProject($id: ID!) {
      deleteProject(id: $id)
    }
  `);
};
