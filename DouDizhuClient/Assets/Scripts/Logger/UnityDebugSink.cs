using Serilog.Core;
using Serilog.Events;
using UnityEngine;

namespace Logger
{
    public class UnityDebugSink : ILogEventSink
    {
        public void Emit(LogEvent logEvent)
        {
            var message = logEvent.RenderMessage();
            switch (logEvent.Level)
            {
                case LogEventLevel.Verbose:
                case LogEventLevel.Debug:
                case LogEventLevel.Information:
                    Debug.Log(message);
                    break;
                case LogEventLevel.Warning:
                    Debug.LogWarning(message);
                    break;
                case LogEventLevel.Error:
                case LogEventLevel.Fatal:
                    Debug.LogError(message);
                    break;
            }
        }
    }
}