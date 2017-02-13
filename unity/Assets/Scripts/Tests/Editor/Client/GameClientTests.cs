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
using Client;
using NSubstitute;
using NUnit.Framework;

namespace Tests.Editor.Client
{
    [TestFixture]
    public class GameClientTests
    {
        private IUnityClient unityServer;

        [SetUp]
        public void Setup()
        {
            unityServer = Substitute.For<IUnityClient>();
        }

        [TearDown]
        public void Teardown()
        {
            GameClient.Stop();
        }

        [Test]
        public void Start()
        {
            GameClient.Start(unityServer, new string[0]);
            unityServer.Received(1).StartClient();
            Assert.Throws<Exception>(() => GameClient.Start(unityServer, new string[0]));
        }

        [Test]
        public void StartWithHost()
        {
            var host = "10.10.10.10";
            GameClient.Start(unityServer, new[] {"-host", host});
            unityServer.Received(1).StartClient();
            unityServer.Received(1).SetHost(host);
            unityServer.Received(1).SetPort(7777);
        }

        [Test]
        public void StartWithPort()
        {
            var args = new[] {"-port", "8080"};
            GameClient.Start(unityServer, args);
            unityServer.Received(1).StartClient();
            unityServer.Received(1).SetHost("localhost");
            unityServer.Received(1).SetPort(8080);
        }

        [Test]
        public void StartWithHostAndPort()
        {
            var host = "10.10.10.10";
            var args = new[] {"-host", host, "-port", "8080"};
            GameClient.Start(unityServer, args);
            unityServer.Received(1).StartClient();
            unityServer.Received(1).SetHost(host);
            unityServer.Received(1).SetPort(8080);
        }
    }
}