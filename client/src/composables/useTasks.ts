import { useMutation } from '@vue/apollo-composable';
import gql from 'graphql-tag';

export const useCreateTaskMutation = () => {
  return useMutation(gql`
    mutation createTask($title: String!, $projectID: ID!) {
      createTask(title: $title, projectID: $projectID) {
        id
        title
        status
      }
    }
  `);
};

export const useUpdateTaskMutation = () => {
  return useMutation(gql`
    mutation updateTask($id: ID!, $title: String, $status: String) {
      updateTask(id: $id, title: $title, status: $status) {
        id
        title
        status
      }
    }
  `);
};

export const useDeleteTaskMutation = () => {
  return useMutation(gql`
    mutation deleteTask($id: ID!) {
      deleteTask(id: $id)
    }
  `);
};
