/**
 * This file provided by Facebook is for non-commercial testing and evaluation
 * purposes only.  Facebook reserves all rights not expressly granted.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * FACEBOOK BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

import {
  commitMutation,
  graphql,
} from 'react-relay';
import { Environment } from 'relay-runtime';

import { DataConstructor } from './typesUtils';

import { TodoList_viewer } from '../__generated__/TodoList_viewer.graphql';
import { MarkAllTodosMutationResponse } from "../__generated__/MarkAllTodosMutation.graphql";

const mutation = graphql`
  mutation MarkAllTodosMutation($input: MarkAllTodosInput!) {
    markAllTodos(input: $input) {
      changedTodos {
        id
        complete
      }
      viewer {
        id
        completedCount
      }
    }
  }
`;

function getOptimisticResponse(complete: boolean, todos: TodoList_viewer["todos"], user: TodoList_viewer): MarkAllTodosMutationResponse {

  type Resp = MarkAllTodosMutationResponse;

  // Developer's duty to return a valid response
  const resp: DataConstructor<Resp> = {};
  resp.markAllTodos = {};
  resp.markAllTodos.viewer = { id: user.id };

  if (todos && todos.edges) {
    resp.markAllTodos.changedTodos = todos.edges
      .filter(edge => edge!.node!.complete !== complete)
      .map(edge => ({
        complete: complete,
        id: edge!.node!.id,
      }));
  }
  if (user.totalCount != null) {
    resp.markAllTodos.viewer.completedCount = complete ?
      user.totalCount :
      0;
  }

  return resp as Resp;
}

function commit(
  environment: Environment,
  complete: boolean,
  todos: TodoList_viewer["todos"],
  user: TodoList_viewer,
) {
  return commitMutation(
    environment,
    {
      mutation,
      variables: {
        input: {complete},
      },
      optimisticResponse: getOptimisticResponse(complete, todos, user),
    }
  );
}

export default {commit};
