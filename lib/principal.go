package lib

type Principal struct {
	userId          uint
	userName        string
	rawToken        string
	token           interface{}
	tokenType       string
	isAuthenticated bool
	privileges      map[string]struct{}
	roles           map[string]struct{}
}

var (
	guest = NewPrincipal("__guest", "None", "None", 0, nil, nil, nil)
)

func NewSimplePrincipal(userName string, userID uint) *Principal {
	return &Principal{
		userId:   userID,
		userName: userName,
	}
}

func NewPrincipal(userName, rawToken, tokenType string, userId uint, privilege, role map[string]struct{}, token interface{}) *Principal {
	return &Principal{
		userId:     userId,
		userName:   userName,
		rawToken:   rawToken,
		token:      token,
		tokenType:  tokenType,
		privileges: privilege,
		roles:      role,
	}
}

func NewAuthenticatedPrincipal(userName, rawToken, tokenType string, userId uint, privilege, role map[string]struct{}, token interface{}) *Principal {
	return &Principal{
		userId:          userId,
		userName:        userName,
		rawToken:        rawToken,
		token:           token,
		tokenType:       tokenType,
		privileges:      privilege,
		roles:           role,
		isAuthenticated: true,
	}
}

func GuestPrincipal() *Principal {
	return guest
}

func (p *Principal) GetUserName() string {
	return p.userName
}

func (p *Principal) GetToken() string {
	return p.rawToken
}

func (p *Principal) GetUserId() uint {
	return p.userId
}

func (p *Principal) GetTokenType() string {
	return p.tokenType
}
func (p *Principal) IsAuthenticated() bool {
	return p.isAuthenticated
}

func (p *Principal) HasPrivilege(privilege string) bool {
	_, exists := p.privileges[privilege]
	return exists
}

func (p *Principal) HasRole(role string) bool {
	_, exists := p.roles[role]
	return exists
}

func (p *Principal) OrPrivilegeAuth(requiredPermissions ...string) bool {
	if requiredPermissions == nil {
		return false
	}
	for _, permission := range requiredPermissions {
		if p.HasPrivilege(permission) {
			return true
		}
	}
	return false
}

func (p *Principal) AndPrivilegeAuth(requiredPermissions ...string) bool {
	if requiredPermissions == nil {
		return false
	}

	permited := 0
	for _, permission := range requiredPermissions {
		if p.HasPrivilege(permission) {
			permited += 1
		}

	}
	return permited == len(requiredPermissions)
}

func (p *Principal) OrRoleAuth(requiredPermissions ...string) bool {
	if requiredPermissions == nil {
		return false
	}
	for _, permission := range requiredPermissions {
		if p.HasRole(permission) {
			return true
		}
	}
	return false
}

func (p *Principal) AndRoleAuth(requiredPermissions ...string) bool {
	if requiredPermissions == nil {
		return false
	}

	permited := 0
	for _, permission := range requiredPermissions {
		if p.HasRole(permission) {
			permited += 1
		}

	}
	return permited == len(requiredPermissions)
}
