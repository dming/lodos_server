1.Client should call "HD_Login" to make his session bind userId and set token.
BTW, client get token by getting the RPC return.

2.When the client logined, he should carry username and token whild sending mqtt publish package.
server would check if he is authenticated bt using the rpc function "Authentication"

3.If authentication fail, client should call "HD_Login" again.
By this way, client can update the session information and get new token.