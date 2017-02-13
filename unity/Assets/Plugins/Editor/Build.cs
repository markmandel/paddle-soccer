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

﻿using UnityEditor;
using UnityEngine;

namespace Plugins.Editor
{
    public class Build
    {
        [MenuItem("Build/Build Game and Server &x")]
        public static void BuildGameAndServer()
        {
            Debug.Log("Building the game client...");
            var originalName = PlayerSettings.productName;

            try
            {
                PlayerSettings.productName = originalName + "-client";
                var options = new BuildPlayerOptions();
                options.scenes = new[] {"Assets/Scenes/SoccerField.unity"};
                options.locationPathName = "Builds/Linux/Game.x86_64";
                options.target = BuildTarget.StandaloneLinux64;
                options.options = BuildOptions.None;
                BuildPipeline.BuildPlayer(options);

                Debug.Log("Building server...");

                PlayerSettings.productName = originalName + "-server";
                options.locationPathName = "Builds/Linux/Server.x86_64";
                options.options = BuildOptions.EnableHeadlessMode;
                BuildPipeline.BuildPlayer(options);
                Debug.Log("Building Complete!");
            }
            finally
            {
                PlayerSettings.productName = originalName;
            }
        }
    }
}
