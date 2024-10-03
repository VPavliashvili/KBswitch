package logging

type Logger interface{
    LogInfo(msg string )
    LogTrace(msg string)
    LogError(msg string)
}
