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

using Game;
using NUnit.Framework;
using UnityEngine;

namespace Tests.Editor.Game
{
    public class GoalsTests
    {
        [Test]
        public void OnBallGoal()
        {
            var check = false;
            var action = Goals.OnBallGoal(_ => check = true);

            var go = new GameObject("Foo");
            var collider = go.AddComponent<SphereCollider>();
            action(collider);
            Assert.IsFalse(check);

            collider.name = Ball.Name;
            action(collider);
            Assert.IsTrue(check);

            Object.DestroyImmediate(go);
        }
    }
}