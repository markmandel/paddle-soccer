using System;
using System.Collections.Generic;
using Network;
using NSubstitute;
using NUnit.Framework;
using Server;
using UnityEngine.Networking;

namespace Tests.Editor.Server
{
    [TestFixture]
    public class GameServerTests
    {
        private IUnityServer unityServer;

        [SetUp]
        public void Setup()
        {
            unityServer = Substitute.For<IUnityServer>();
            unityServer.StartServer().Returns(true);
            GameServer.Start(unityServer);
        }

        [TearDown]
        public void Teardown()
        {
            GameServer.Stop();
        }

        [Test]
        public void Start()
        {
            // reset state
            GameServer.Stop();
            unityServer = Substitute.For<IUnityServer>();

            unityServer.StartServer().Returns(false);
            Assert.Throws<Exception>(() => GameServer.Start(unityServer));
            unityServer.StartServer().Returns(true);
            GameServer.Start(unityServer);
            unityServer.Received(2).StartServer();
            Assert.Throws<Exception>(() => GameServer.Start(unityServer));
        }

        [Test]
        public void OnServerConnect()
        {
            var conn = Substitute.For<NetworkConnection>();

            GameServer.OnServerConnect(conn);
            conn.Received(0).Disconnect();
            GameServer.OnServerConnect(conn);
            conn.Received(0).Disconnect();
            GameServer.OnServerConnect(conn);
            conn.Received(1).Disconnect();
        }

        [Test]
        public void OnServerAddPlayer()
        {
            var isReady = false;
            GameServer.OnGameReady += () => isReady = true;
            Assert.False(isReady);

            var fixtures = new List<AddPlayerFixure>
            {
                new AddPlayerFixure(0, false),
                new AddPlayerFixure(1, true)
            };

            fixtures.ForEach(x =>
            {
                GameServer.OnServerAddPlayer(x.playerCount);
                Assert.AreEqual(x.isReady, isReady);
            });
        }

        private class AddPlayerFixure
        {
            public int playerCount;
            public bool isReady;

            public AddPlayerFixure(int playerCount, bool isReady)
            {
                this.playerCount = playerCount;
                this.isReady = isReady;
            }
        }
    }
}