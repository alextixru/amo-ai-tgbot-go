// Copyright 2025 Google LLC
// SPDX-License-Identifier: Apache-2.0

package geminicli

import (
	"github.com/firebase/genkit/go/ai"
)

var (
	// Multimodal describes model capabilities for multimodal Gemini models.
	Multimodal = ai.ModelSupports{
		Multiturn:   true,
		Tools:       true,
		ToolChoice:  true,
		SystemRole:  true,
		Media:       true,
		Constrained: ai.ConstrainedSupportNoTools,
	}

	// Gemini3Support includes reasoning (thought) capabilities.
	Gemini3Support = ai.ModelSupports{
		Multiturn:   true,
		Tools:       true,
		ToolChoice:  true,
		SystemRole:  true,
		Media:       true,
		Constrained: ai.ConstrainedSupportNoTools,
		// Reasoning: true, // Note: Genkit Go ai.ModelSupports doesn't have a Reasoning field yet, but we handle it via parts.
	}
)

const (
	// PA-API Specific Models
	Gemini3ProPreview   = "gemini-3-pro-preview"
	Gemini3FlashPreview = "gemini-3-flash-preview"
	Gemini25Pro         = "gemini-2.5-pro"
	Gemini25Flash       = "gemini-2.5-flash"
	Gemini25FlashLite   = "gemini-2.5-flash-lite"

	// Stable references for transition
)

// SupportedModels defines the models used by the Code Assist provider.
var SupportedModels = map[string]ai.ModelOptions{
	Gemini3ProPreview: {
		Label:    "Gemini 3 Pro Preview",
		Supports: &Gemini3Support,
		Stage:    ai.ModelStageUnstable,
	},
	Gemini3FlashPreview: {
		Label:    "Gemini 3 Flash Preview",
		Supports: &Gemini3Support,
		Stage:    ai.ModelStageUnstable,
	},
	Gemini25Pro: {
		Label:    "Gemini 2.5 Pro",
		Supports: &Multimodal,
		Stage:    ai.ModelStageStable,
	},
	Gemini25Flash: {
		Label:    "Gemini 2.5 Flash",
		Supports: &Multimodal,
		Stage:    ai.ModelStageStable,
	},
	Gemini25FlashLite: {
		Label:    "Gemini 2.5 Flash Lite",
		Supports: &Multimodal,
		Stage:    ai.ModelStageStable,
	},
}
