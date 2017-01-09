using UnityEngine;
using UnityEngine.Networking;
using Server;

namespace Network
{
    public class GameNetwork : NetworkManager, IUnityServer
    {
        public readonly string Version = "0.2";

        /// <summary>
        /// How many players have joined the game?
        /// </summary>
        private int connCount;

        /// <summary>
        /// Starts eiter a client or a server - depending on if headless or not.
        /// </summary>
        private void Start()
        {
            Debug.LogFormat("[GameNetwork] Starting Client or Server? {0}", Version);

            if (PlayerInfo.IsHeadless())
            {
                Debug.Log("[GameNetwork] Starting Server");
                GameServer.Start(this);
            }
            else
            {
                Debug.Log("[GameNetwork] Starting Client");
                StartClient();
            }
        }

        // --- Server Commands ---
        public override void OnServerConnect(NetworkConnection conn)
        {
            base.OnServerConnect(conn);
            GameServer.OnServerConnect(conn);
        }

        public override void OnStopServer()
        {
            base.OnStopServer();
            GameServer.Stop();
        }

        public override void OnServerAddPlayer(NetworkConnection conn, short playerControllerId)
        {
            base.OnServerAddPlayer(conn, playerControllerId);
            GameServer.OnServerAddPlayer(numPlayers);
        }

        // --- Client Commands ---
    }
}