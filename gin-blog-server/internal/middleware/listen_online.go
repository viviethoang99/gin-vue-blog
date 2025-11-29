package middleware

import (
	"context"
	"fmt"
	g "gin-blog/internal/global"
	"gin-blog/internal/handle"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

// Listen online status middleware
// On login: remove user's force-offline flag
// On logout: add user's online flag
func ListenOnline() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		rdb := c.MustGet(g.CTX_RDB).(*redis.Client)

		auth, err := handle.CurrentUserAuth(c)
		if err != nil {
			handle.ReturnError(c, g.ErrUserAuth, err)
			return
		}

		onlineKey := g.ONLINE_USER + strconv.Itoa(auth.ID)
		offlineKey := g.OFFLINE_USER + strconv.Itoa(auth.ID)

		// Check whether current user is forced offline
		if rdb.Exists(ctx, offlineKey).Val() == 1 {
			fmt.Println("User was forced offline")
			handle.ReturnError(c, g.ErrForceOffline, nil)
			c.Abort()
			return
		}

		// Each request refreshes online status in Redis: reset to 10 minutes
		rdb.Set(ctx, onlineKey, auth, 10*time.Minute)
		c.Next()
	}
}
