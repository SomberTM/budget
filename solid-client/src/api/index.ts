import Axios, { AxiosHeaders, HttpStatusCode } from "axios";
import { ApiResponse } from "./response";
import * as Users from "./users";

export namespace Api {
  export const axios = Axios.create({
    baseURL: "/api/",
    transformResponse: [
      (data) => {
        if (!data) return;
        return JSON.parse(data);
      },
      function (
        data: any,
        headers: AxiosHeaders,
        status?: number
      ): ApiResponse {
        if (!status) return data;

        if (status >= 200 && status < 300) {
          return { success: true, data };
        }

        return {
          success: false,
          detailedError: data,
          status,
          statusText: HttpStatusCode[status],
        };
      },
    ],
    validateStatus: () => true,
  });

  /**
   * - Shorthand for `Api.axios.get`
   * - Use this when you just want the ApiResponse out and couldn't care about headers or other http related stuff
   * @returns
   */
  export const get = <T>(
    ...args: Parameters<typeof axios.get>
  ): Promise<ApiResponse<T>> =>
    axios.get(...args).then((response) => response.data);

  export const post = <T>(
    ...args: Parameters<typeof axios.post>
  ): Promise<ApiResponse<T>> =>
    axios.post(...args).then((response) => response.data);

  export const users = Users;
}
