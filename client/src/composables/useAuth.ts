import { useMutation } from '@vue/apollo-composable';
import gql from 'graphql-tag';

export const useLogin = () => {
  return useMutation(gql`
    mutation login($username: String!, $password: String!) {
      login(username: $username, password: $password) {
        token
        user {
          id
          username
          email
        }
      }
    }
  `);
};

export const useRegister = () => {
  return useMutation(gql`
    mutation createUser($username: String!, $email: String!, $password: String!) {
      createUser(username: $username, email: $email, password: $password) {
        id
        username
        email
      }
    }
  `);
};
