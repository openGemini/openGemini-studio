// Copyright 2026 openGemini Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

"use client"

import { ref, watch } from "vue"

export type Theme = "light" | "dark"

let themeInstance: ReturnType<typeof createTheme> | null = null

function createTheme() {
  const theme = ref<Theme>("light")
  const isInitialized = ref(false)

  // 应用主题到 document
  const applyTheme = (newTheme: Theme) => {
    if (typeof document !== "undefined") {
      document.documentElement.setAttribute("data-theme", newTheme)
    }
  }

  // 初始化主题
  const initTheme = () => {
    if (typeof localStorage !== "undefined" && !isInitialized.value) {
      const savedTheme = localStorage.getItem("theme") as Theme | null
      if (savedTheme) {
        theme.value = savedTheme
      }
      applyTheme(theme.value)
      isInitialized.value = true
    }
  }

  // 立即初始化
  initTheme()

  const toggleTheme = () => {
    theme.value = theme.value === "dark" ? "light" : "dark"
  }

  watch(theme, (newTheme) => {
    if (typeof localStorage !== "undefined") {
      localStorage.setItem("theme", newTheme)
    }
    applyTheme(newTheme)
  })

  return {
    theme,
    toggleTheme,
    initTheme,
  }
}

export function useTheme() {
  if (!themeInstance) {
    themeInstance = createTheme()
  }
  return themeInstance
}
