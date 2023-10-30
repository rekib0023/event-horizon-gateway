package utils

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

func GetHttpStatusCode(code codes.Code) int {
	switch code {
	case 0: // Status.OK
		return http.StatusOK
	case 1: // Status.CANCELLED
		return http.StatusConflict
	case 2: // Status.UNKNOWN
		return http.StatusInternalServerError
	case 3: // Status.INVALID_ARGUMENT
		return http.StatusBadRequest
	case 4: // Status.DEADLINE_EXCEEDED
		return http.StatusRequestTimeout
	case 5: // Status.NOT_FOUND
		return http.StatusNotFound
	case 6: // Status.ALREADY_EXISTS
		return http.StatusConflict
	case 7: // Status.PERMISSION_DENIED
		return http.StatusForbidden
	case 8: // Status.RESOURCE_EXHAUSTED
		return http.StatusTooManyRequests
	case 9: // Status.FAILED_PRECONDITION
		return http.StatusPreconditionFailed
	case 10: // Status.ABORTED
		return http.StatusConflict
	case 11: // Status.OUT_OF_RANGE
		return http.StatusBadRequest
	case 12: // Status.UNIMPLEMENTED
		return http.StatusNotImplemented
	case 13: // Status.INTERNAL
		return http.StatusInternalServerError
	case 14: // Status.UNAVAILABLE
		return http.StatusServiceUnavailable
	case 15: // Status.DATA_LOSS
		return http.StatusInternalServerError
	case 16: // Status.UNAUTHENTICATED
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
