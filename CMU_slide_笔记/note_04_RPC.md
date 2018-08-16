**RPC挑战**
1.lost message or server,client crash
2.语义
exactly once :不可能实现
at least once :仅仅应用于幂等操作
at most once :需要思考例外情况
例如: client获取server lock,at most once并不能保证会有replay   =》 handle machiune failures =》需要额外的查询自身state的require
