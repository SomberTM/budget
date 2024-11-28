type OkApiResponse<T = unknown> = {
  success: true;
  data: T;
};

type ApiError = {
  message: string;
};

type ErrorApiResponse = {
  success: false;
  detailedError?: ApiError;
  status: number;
  statusText: string;
};

export type ApiResponse<T = unknown> = OkApiResponse<T> | ErrorApiResponse;
