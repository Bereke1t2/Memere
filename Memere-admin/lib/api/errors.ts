export class ApiError extends Error {
  code: string;
  status: number;
  details?: unknown;

  constructor(code: string, message: string, status: number, details?: unknown) {
    super(message);
    this.name = "ApiError";
    this.code = code;
    this.status = status;
    this.details = details;
  }
}

const CODE_MESSAGES: Record<string, string> = {
  INVALID_CREDENTIALS: "Email or password is incorrect.",
  RESOURCE_NOT_FOUND: "The requested resource was not found.",
  NOT_FOUND: "The requested resource was not found.",
  FORBIDDEN: "You do not have permission to perform this action.",
  UNAUTHORIZED: "Your session has expired. Please log in again.",
  SESSION_EXPIRED: "Your session has expired. Please log in again.",
  IDEMPOTENCY_KEY_REQUIRED: "A unique request key is required.",
  VALIDATION_ERROR: "One or more fields are invalid.",
  CONFLICT: "This action conflicts with an existing resource.",
  VIDEO_EXISTS: "A video already exists for this lesson. Delete the existing video or create a new lesson.",
  INVALID_ID: "The provided ID is invalid.",
  INTERNAL_SERVER_ERROR: "An unexpected error occurred. Please try again.",
};

export function friendlyMessage(err: ApiError): string {
  if (err.code === "VIDEO_EXISTS" || (err.code === "CONFLICT" && (err.message?.toLowerCase().includes("video") || err.message?.toLowerCase().includes("conflict")))) {
    return "A video already exists for this lesson. Please delete the existing video or create a new lesson.";
  }
  return CODE_MESSAGES[err.code] ?? err.message;
}
