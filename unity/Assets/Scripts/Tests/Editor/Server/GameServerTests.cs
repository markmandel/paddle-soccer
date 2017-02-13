// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

﻿using System;
using System.Collections.Generic;
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
            unityServer.Received(10).StartServer();

            unityServer.ClearReceivedCalls();
            unityServer.StartServer().Returns(true);
            GameServer.Start(unityServer);
            unityServer.Received(1).StartServer();
            Assert.Throws<Exception>(() => GameServer.Start(unityServer));
        }

        [Test]
        public void SelectPort()
        {
            Environment.SetEnvironmentVariable(GameServer.MinPortEnv, null);
            Environment.SetEnvironmentVariable(GameServer.MaxPortEnv, null);

            for (var i = 0; i < 100; i++)
            {
                GameServer.Stop();
                unityServer.ClearReceivedCalls();
                GameServer.Start(unityServer);
                unityServer.Received(1).SetPort(Arg.Is<int>(x => 7000 <= x && x <= 8000));
            }

            Environment.SetEnvironmentVariable(GameServer.MinPortEnv, "10");
            Environment.SetEnvironmentVariable(GameServer.MaxPortEnv, "100");

            for (var i = 0; i < 100; i++)
            {
                GameServer.Stop();
                unityServer.ClearReceivedCalls();
                GameServer.Start(unityServer);
                unityServer.Received(1).SetPort(Arg.Is<int>(x => 10 <= x && x <= 100));
            }
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
                new AddPlayerFixure(1, false),
                new AddPlayerFixure(2, true)
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