package registration

import "errors"

var RegistrationFailedError = errors.New(RegistrationFailedMessage)
var AccountAlreadyRegisteredError = errors.New(AccountAlreadyRegisteredErrorMessage)
