package api

import (
	"fmt"
	"time"

	"google.golang.org/grpc/metadata"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ellcrys/util"
	"github.com/ncodes/cocoon/core/api/api/proto"
	orderer_proto "github.com/ncodes/cocoon/core/orderer/proto"
	"github.com/ncodes/cocoon/core/types"
	"golang.org/x/crypto/bcrypt"
	context "golang.org/x/net/context"
)

// authenticateToken validates a token
func (api *API) authenticateToken(tokenStr string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		tokenType := token.Claims.(jwt.MapClaims)["type"]
		if tokenType == "token.cli" {
			return []byte(util.Env("API_SIGN_KEY", "secret")), nil
		}
		return nil, fmt.Errorf("unknown token type")
	})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	if err = claims.Valid(); err != nil {
		return nil, err
	}

	return claims, nil
}

// checkCtxAccessToken authenticates the access token in the context
func (api *API) checkCtxAccessToken(ctx context.Context) (jwt.MapClaims, error) {
	md, _ := metadata.FromContext(ctx)
	accessTokens := md["access_token"]
	if accessTokens == nil {
		return nil, fmt.Errorf("access token is required")
	}
	return api.authenticateToken(accessTokens[0])
}

// makeAuthToken creates a session token
func makeAuthToken(id, identity, _type string, exp int64, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"id":       id,
		"identity": identity,
		"type":     _type,
		"exp":      exp,
	})
	return token.SignedString([]byte(secret))
}

// Login authenticates a user and returns a JWT token
func (api *API) Login(ctx context.Context, req *proto.LoginRequest) (*proto.Response, error) {

	ordererConn, err := api.ordererDiscovery.GetGRPConn()
	if err != nil {
		return nil, err
	}
	defer ordererConn.Close()

	resp, err := api.GetIdentity(ctx, &proto.GetIdentityRequest{
		Email: req.GetEmail(),
	})

	if err != nil && err != types.ErrIdentityNotFound {
		return nil, err
	} else if err == types.ErrIdentityNotFound {
		return nil, fmt.Errorf("email or password are invalid")
	}

	var identity types.Identity
	err = util.FromJSON(resp.GetBody(), &identity)
	if err != nil {
		return nil, fmt.Errorf("malformed identity")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(identity.Password), []byte(req.GetPassword())); err != nil {
		return nil, fmt.Errorf("email or password are invalid")
	}

	sessionID := util.Sha256(util.UUID4())
	identity.ClientSessions = append(identity.ClientSessions, sessionID)
	key := util.Env("API_SIGN_KEY", "secret")
	ss, err := makeAuthToken(sessionID, identity.GetID(), "token.cli", time.Now().AddDate(0, 1, 0).Unix(), key)
	if err != nil {
		apiLog.Error(err.Error())
		return nil, fmt.Errorf("failed to create session token")
	}

	// overwrite identity to store the newly created client session.
	odc := orderer_proto.NewOrdererClient(ordererConn)
	_, err = odc.Put(ctx, &orderer_proto.PutTransactionParams{
		CocoonID:   "",
		LedgerName: types.GetGlobalLedgerName(),
		Transactions: []*orderer_proto.Transaction{
			&orderer_proto.Transaction{
				Id:        util.UUID4(),
				Key:       api.makeIdentityKey(identity.Email),
				Value:     string(identity.ToJSON()),
				CreatedAt: time.Now().Unix(),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &proto.Response{
		Status: 200,
		Body:   []byte(ss),
	}, nil
}
