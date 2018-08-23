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

import { Todo_todo } from '../__generated__/Todo_todo.graphql';
import { Todo_viewer } from '../__generated__/Todo_viewer.graphql';
import { ChangeTodoStatusMutationResponse } from '../__generated__/ChangeTodoStatusMutation.graphql';

const mutation = graphql`
  mutation ChangeTodoStatusMutation($input: ChangeTodoStatusInput!) {
    changeTodoStatus(input: $input) {
      todo {
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

function getOptimisticResponse(complete: boolean, todo: Todo_todo, user: Todo_viewer): ChangeTodoStatusMutationResponse {

  type Resp = ChangeTodoStatusMutationResponse;
  // Developer's duty to return a valid response
  const resp: DataConstructor<Resp> = {};
  resp.changeTodoStatus = {};
  resp.changeTodoStatus.viewer = {id: user.id};

  if (user.completedCount != null) {
    resp.changeTodoStatus.viewer.completedCount = complete ?
      user.completedCount + 1 :
      user.completedCount - 1;
  }
  resp.changeTodoStatus.todo = {
    complete: complete,
    id: todo.id,
  };

  return resp as Resp;
}

function commit(
  environment: Environment,
  complete: boolean,
  todo: Todo_todo,
  user: Todo_viewer,
) {
  return commitMutation(
    environment,
    {
      mutation,
      variables: {
        input: {complete, id: todo.id},
      },
      optimisticResponse: getOptimisticResponse(complete, todo, user),
    }
  );
}

export default {commit};
