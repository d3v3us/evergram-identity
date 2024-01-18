package auth

import (
	"errors"
)

var ErrInvalidCredentials = errors.New(InvalidCredentialsMessage)
var ErrAuthentication = errors.New(AuthenticationErrorMessage)
var ErrUserWithSpecifiedCredentialsNotFound = errors.New(UserWithSpecifiedCredentialsNotFoundMessage)
var ErrInvalidJwtToken = errors.New(InvalidJwtTokenMessage)
var ErrFaildToExtractJwtToken = errors.New(FaildToExtractJwtTokenMessage)
var ErrInvalidJwtTokenSigningMethod = errors.New(InvalidJwtTokenSigningMethodMessage)
var ErrJwtTokenAuthenticationFailed = errors.New(JwtTokenAuthenticationFailedMessage)

var ErrInvalidIssuer = errors.New(InvalidIssuerMessage)
var ErrInvalidSubject = errors.New(InvalidSubjectMessage)
var ErrInvalidID = errors.New(InvalidIDMessage)
var ErrInvalidAudience = errors.New(InvalidAudienceMessage)
var ErrMissingOrMailformedJWT = errors.New(MissingOrMailformedJWTMessage)
