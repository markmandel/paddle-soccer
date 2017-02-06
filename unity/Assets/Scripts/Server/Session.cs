using System;

namespace Server
{
    /// <summary>
    /// Serialisable session - for talking to the session server
    /// </summary>
    [Serializable]
    public class Session
    {
        public string id;
        public int port;
    }
}