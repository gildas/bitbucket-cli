package pipeline

import (
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

// Target represents the target of a pipeline (branch, tag, etc.)
type Target interface {
	core.TypeCarrier
	GetDestination() string
	GetCommit() commit.CommitReference
}

var targetRegistry = core.TypeRegistry{}

// UnmarshalTarget unmarshals a Target from JSON data
func UnmarshalTarget(payload []byte) (Target, error) {
	target, err := targetRegistry.UnmarshalJSON(payload)
	if err != nil {
		errStr := err.Error()
		
		if strings.HasPrefix(errStr, "Missing JSON Property") {
			return nil, errors.JSONUnmarshalError.Wrap(errors.ArgumentMissing.With("type"))
		}
		
		// Check for "Unsupported Type" error and extract the type name
		if strings.HasPrefix(errStr, "Unsupported Type ") {
			// Extract the type name from the error message: 'Unsupported Type "typename"'
			const prefix = "Unsupported Type "
			remainder := errStr[len(prefix):]
			// Remove surrounding quotes
			typeName := strings.Trim(remainder, `"`)
			
			// Get list of supported types
			supportedTypes := targetRegistry.SupportedTypes()
			
			return nil, errors.JSONUnmarshalError.Wrap(
				errors.InvalidType.With(typeName, strings.Join(supportedTypes, ", ")),
			)
		}
		
		return nil, err
	}
	return target.(Target), nil
}
