# AI Features Implementation Proposal

## Overview

This document outlines a comprehensive plan to add AI capabilities to Canvas CLI, transforming it from a data access tool into an intelligent assistant for educators.

## Phase 1: Foundation (Week 1-2)

### 1.1 AI Provider Abstraction

Create `/internal/ai/client.go`:

```go
package ai

import "context"

type Provider interface {
    Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)
    Stream(ctx context.Context, req *CompletionRequest) (<-chan CompletionChunk, error)
}

type CompletionRequest struct {
    Prompt       string
    SystemPrompt string
    MaxTokens    int
    Temperature  float64
    Model        string
}

type CompletionResponse struct {
    Content      string
    TokensUsed   int
    FinishReason string
}

// Factory function
func NewProvider(cfg *config.AIConfig) (Provider, error) {
    switch cfg.Provider {
    case "anthropic":
        return NewAnthropicProvider(cfg.APIKey, cfg.Model)
    case "openai":
        return NewOpenAIProvider(cfg.APIKey, cfg.Model)
    case "local":
        return NewLocalProvider(cfg.Endpoint)
    default:
        return nil, fmt.Errorf("unknown AI provider: %s", cfg.Provider)
    }
}
```

### 1.2 Configuration Schema

Extend `/internal/config/config.go`:

```go
type AIConfig struct {
    Enabled         bool    `mapstructure:"enabled"`
    Provider        string  `mapstructure:"provider"` // anthropic, openai, local
    APIKey          string  `mapstructure:"api_key"`
    Model           string  `mapstructure:"model"`
    MaxTokens       int     `mapstructure:"max_tokens"`
    Temperature     float64 `mapstructure:"temperature"`
    CacheResponses  bool    `mapstructure:"cache_responses"`
    AnonymizeData   bool    `mapstructure:"anonymize_data"`
    LocalEndpoint   string  `mapstructure:"local_endpoint"` // for ollama
}
```

### 1.3 Command Infrastructure

Create `/commands/ai.go`:

```go
package commands

import (
    "github.com/jjuanrivvera/canvas-cli/commands/internal/options"
    "github.com/spf13/cobra"
)

func NewAICmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "ai",
        Short: "AI-powered features for Canvas",
        Long:  `Intelligent analysis and automation using AI`,
    }

    cmd.AddCommand(
        newAISuggestGradeCmd(),
        newAIAnalyzeCourseCmd(),
        newAISuggestQuestionsCmd(),
        newAIGenerateQuizCmd(),
        newAIQueryCmd(),
    )

    return cmd
}
```

## Phase 2: Grading Assistant (Week 3)

### High-Value Feature: Intelligent Grade Suggestions

Create `/internal/ai/grading.go`:

```go
package ai

import (
    "context"
    "fmt"
    "strings"

    "github.com/jjuanrivvera/canvas-cli/internal/api"
)

type GradingService struct {
    provider     Provider
    apiClient    *api.Client
    anonymize    bool
}

type GradeSuggestion struct {
    SuggestedScore float64
    MaxPoints      float64
    Explanation    string
    RubricScores   map[string]float64 // criterion_id -> score
    Confidence     float64
    Warnings       []string
}

func (s *GradingService) SuggestGrade(
    ctx context.Context,
    assignmentID, submissionID int64,
) (*GradeSuggestion, error) {
    // 1. Fetch assignment with rubric
    assignment, err := s.apiClient.Assignments.Get(ctx, assignmentID, []string{"rubric"})
    if err != nil {
        return nil, fmt.Errorf("fetching assignment: %w", err)
    }

    // 2. Fetch submission
    submission, err := s.apiClient.Submissions.Get(ctx, assignmentID, submissionID)
    if err != nil {
        return nil, fmt.Errorf("fetching submission: %w", err)
    }

    // 3. Build grading prompt
    prompt := s.buildGradingPrompt(assignment, submission)

    // 4. Call AI provider
    resp, err := s.provider.Complete(ctx, &CompletionRequest{
        Prompt:       prompt,
        SystemPrompt: systemPromptGrading,
        MaxTokens:    2000,
        Temperature:  0.2, // Lower for consistency
    })
    if err != nil {
        return nil, fmt.Errorf("AI completion: %w", err)
    }

    // 5. Parse response
    return s.parseGradeSuggestion(resp.Content, assignment.PointsPossible)
}

func (s *GradingService) buildGradingPrompt(
    assignment *api.Assignment,
    submission *api.Submission,
) string {
    var sb strings.Builder

    sb.WriteString("# Assignment Details\n\n")
    sb.WriteString(fmt.Sprintf("**Title:** %s\n", assignment.Name))
    sb.WriteString(fmt.Sprintf("**Points Possible:** %.2f\n\n", assignment.PointsPossible))

    if assignment.Description != "" {
        sb.WriteString("**Description:**\n")
        sb.WriteString(cleanHTML(assignment.Description))
        sb.WriteString("\n\n")
    }

    // Include rubric if available
    if len(assignment.Rubric) > 0 {
        sb.WriteString("# Rubric Criteria\n\n")
        for _, criterion := range assignment.Rubric {
            sb.WriteString(fmt.Sprintf("## %s (%.2f points)\n", criterion.Description, criterion.Points))
            sb.WriteString(fmt.Sprintf("%s\n\n", criterion.LongDescription))

            // Rating levels
            for _, rating := range criterion.Ratings {
                sb.WriteString(fmt.Sprintf("- **%s** (%.2f pts): %s\n",
                    rating.Description, rating.Points, rating.LongDescription))
            }
            sb.WriteString("\n")
        }
    }

    // Student submission
    sb.WriteString("# Student Submission\n\n")

    if s.anonymize {
        sb.WriteString("**[Student name anonymized]**\n\n")
    } else {
        sb.WriteString(fmt.Sprintf("**Student:** %s\n\n", submission.User.Name))
    }

    if submission.Body != "" {
        sb.WriteString(cleanHTML(submission.Body))
    }

    sb.WriteString("\n\n# Task\n\n")
    sb.WriteString("Analyze this submission and provide:\n")
    sb.WriteString("1. A suggested score (0-" + fmt.Sprintf("%.2f", assignment.PointsPossible) + ")\n")
    sb.WriteString("2. Scores for each rubric criterion (if rubric provided)\n")
    sb.WriteString("3. Detailed explanation of the grade\n")
    sb.WriteString("4. Constructive feedback for the student\n")
    sb.WriteString("5. Confidence level (0-1)\n\n")
    sb.WriteString("Format your response as JSON:\n")
    sb.WriteString(`{
  "score": <number>,
  "rubric_scores": {"criterion_id": <score>, ...},
  "explanation": "<string>",
  "feedback": "<string>",
  "confidence": <number>
}`)

    return sb.String()
}

const systemPromptGrading = `You are an expert educational assessment assistant.
Your role is to evaluate student work fairly and consistently based on provided rubrics and assignment criteria.

Guidelines:
- Be objective and evidence-based in your assessment
- Reference specific rubric criteria in your evaluation
- Provide constructive, actionable feedback
- Acknowledge strong work and areas for improvement
- If the submission is incomplete or off-topic, reflect that in the score
- Express lower confidence when criteria are ambiguous or submission is unclear
- Never inflate grades; be honest about quality`
```

### Command Implementation

Create `/commands/ai_grading.go`:

```go
func newAISuggestGradeCmd() *cobra.Command {
    opts := &options.AISuggestGradeOptions{}

    cmd := &cobra.Command{
        Use:   "suggest-grade",
        Short: "AI-powered grade suggestion for a submission",
        Long: `Analyzes a submission using AI and suggests a grade based on
the assignment rubric and description.`,
        Example: `  canvas ai suggest-grade \
    --assignment-id 12345 \
    --submission-id 67890 \
    --explain`,
        RunE: func(cmd *cobra.Command, args []string) error {
            logger := logging.NewCommandLogger(globalDebugFlag)

            if err := opts.Validate(); err != nil {
                return err
            }

            client, _ := getAPIClient()
            aiProvider, err := getAIProvider()
            if err != nil {
                return fmt.Errorf("AI not configured: %w\n\nRun: canvas config set ai.provider anthropic", err)
            }

            return runAISuggestGrade(cmd.Context(), client, aiProvider, opts, logger)
        },
    }

    cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
    cmd.Flags().Int64Var(&opts.SubmissionID, "submission-id", 0, "Submission ID (required)")
    cmd.Flags().BoolVar(&opts.Explain, "explain", false, "Show detailed explanation")
    cmd.Flags().BoolVar(&opts.ApplyGrade, "apply", false, "Apply the suggested grade")
    cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Show what would be applied without applying")

    cmd.MarkFlagRequired("assignment-id")
    cmd.MarkFlagRequired("submission-id")

    return cmd
}

func runAISuggestGrade(
    ctx context.Context,
    client *api.Client,
    provider ai.Provider,
    opts *options.AISuggestGradeOptions,
    logger *logging.CommandLogger,
) error {
    logger.LogCommandStart(ctx, "ai.suggest-grade", map[string]interface{}{
        "assignment_id": opts.AssignmentID,
        "submission_id": opts.SubmissionID,
    })

    gradingService := ai.NewGradingService(provider, client, true)

    fmt.Println("🤖 Analyzing submission with AI...")

    suggestion, err := gradingService.SuggestGrade(ctx, opts.AssignmentID, opts.SubmissionID)
    if err != nil {
        return fmt.Errorf("generating suggestion: %w", err)
    }

    // Display results
    fmt.Printf("\n✓ Grade Suggestion\n\n")
    fmt.Printf("Suggested Score: %.2f / %.2f\n", suggestion.SuggestedScore, suggestion.MaxPoints)
    fmt.Printf("Confidence: %.0f%%\n\n", suggestion.Confidence*100)

    if len(suggestion.RubricScores) > 0 {
        fmt.Println("Rubric Breakdown:")
        for criterionID, score := range suggestion.RubricScores {
            fmt.Printf("  - %s: %.2f pts\n", criterionID, score)
        }
        fmt.Println()
    }

    if opts.Explain {
        fmt.Println("Explanation:")
        fmt.Println(wrapText(suggestion.Explanation, 80))
        fmt.Println()
    }

    if len(suggestion.Warnings) > 0 {
        fmt.Println("⚠️  Warnings:")
        for _, warning := range suggestion.Warnings {
            fmt.Printf("  - %s\n", warning)
        }
        fmt.Println()
    }

    if opts.ApplyGrade {
        if !opts.DryRun {
            fmt.Print("Apply this grade? [y/N]: ")
            var response string
            fmt.Scanln(&response)
            if strings.ToLower(response) != "y" {
                fmt.Println("Grade not applied")
                return nil
            }
        }

        if opts.DryRun {
            fmt.Println("🔍 Dry run - would apply:")
        } else {
            fmt.Println("📝 Applying grade...")
        }

        // Apply grade (if not dry run)
        if !opts.DryRun {
            err := client.Submissions.Grade(ctx, opts.AssignmentID, opts.SubmissionID, &api.GradeParams{
                Score:   &suggestion.SuggestedScore,
                Comment: "AI-suggested grade: " + suggestion.Explanation,
            })
            if err != nil {
                return fmt.Errorf("applying grade: %w", err)
            }
            fmt.Println("✓ Grade applied successfully")
        }
    }

    logger.LogCommandComplete(ctx, "ai.suggest-grade", 1)
    return nil
}
```

## Phase 3: Analytics Insights (Week 4)

Create `/internal/ai/analytics.go`:

```go
package ai

type AnalyticsService struct {
    provider  Provider
    apiClient *api.Client
}

type CourseInsights struct {
    Summary            string
    AtRiskStudents     []StudentRisk
    AssignmentIssues   []AssignmentIssue
    EngagementTrends   string
    Recommendations    []string
}

type StudentRisk struct {
    UserID            int64
    Name              string
    RiskLevel         string // high, medium, low
    Reasons           []string
    SuggestedActions  []string
}

func (s *AnalyticsService) AnalyzeCourse(
    ctx context.Context,
    courseID int64,
) (*CourseInsights, error) {
    // 1. Gather analytics data
    studentSummaries, _ := s.apiClient.Analytics.GetStudentSummaries(ctx, courseID)
    courseAnalytics, _ := s.apiClient.Analytics.GetCourseActivity(ctx, courseID)
    assignments, _ := s.apiClient.Assignments.List(ctx, courseID)

    // 2. Build analytics prompt
    prompt := s.buildAnalyticsPrompt(courseID, studentSummaries, courseAnalytics, assignments)

    // 3. Get AI insights
    resp, err := s.provider.Complete(ctx, &CompletionRequest{
        Prompt:       prompt,
        SystemPrompt: systemPromptAnalytics,
        MaxTokens:    3000,
        Temperature:  0.3,
    })
    if err != nil {
        return nil, err
    }

    // 4. Parse and structure insights
    return s.parseInsights(resp.Content, studentSummaries)
}

const systemPromptAnalytics = `You are an expert educational data analyst.
Analyze course analytics data and provide actionable insights for instructors.

Focus on:
- Identifying students who may need additional support
- Detecting concerning patterns in engagement or performance
- Highlighting successful strategies
- Recommending evidence-based interventions

Be specific, actionable, and supportive in tone.`
```

## Phase 4: Discussion & Content Generation (Week 5-6)

### Discussion Enhancement

```go
// /internal/ai/discussion.go

func (s *DiscussionService) SuggestQuestions(
    ctx context.Context,
    topicID int64,
    count int,
) ([]string, error) {
    // Analyze existing discussion
    // Generate thought-provoking follow-up questions
}

func (s *DiscussionService) SummarizeDiscussion(
    ctx context.Context,
    topicID int64,
) (*DiscussionSummary, error) {
    // Extract key themes
    // Identify consensus and disagreement
    // Highlight insightful contributions
}
```

### Content Generation

```go
// /internal/ai/content.go

func (s *ContentService) GenerateQuiz(
    ctx context.Context,
    courseID int64,
    topic string,
    opts *QuizOptions,
) (*api.Quiz, error) {
    // Generate questions based on course materials
    // Create distractors for multiple choice
    // Provide answer explanations
}

func (s *ContentService) SuggestRubric(
    ctx context.Context,
    assignmentID int64,
) (*api.Rubric, error) {
    // Analyze assignment description
    // Suggest evaluation criteria
    // Create rating scales
}
```

## Phase 5: Natural Language Query (Week 7)

Create intelligent query interface:

```bash
canvas ai query "Show me students who haven't submitted assignment 5"
# Translates to:
# canvas submissions list --assignment-id 5 --workflow-state unsubmitted

canvas ai query "Which course has the most overdue assignments?"
# Analyzes across courses and reports findings
```

```go
// /internal/ai/query.go

func (s *QueryService) ExecuteNaturalLanguageQuery(
    ctx context.Context,
    query string,
) (*QueryResult, error) {
    // 1. Parse intent
    intent, err := s.parseQuery(ctx, query)

    // 2. Map to Canvas API calls
    apiCalls := s.mapToAPICalls(intent)

    // 3. Execute calls
    results, err := s.executeCalls(ctx, apiCalls)

    // 4. Format response
    return s.formatResults(ctx, results, query)
}
```

## Configuration & Setup

### Initial Setup Command

```bash
canvas ai setup
```

Interactive wizard:
1. Choose provider (Anthropic/OpenAI/Local)
2. Enter API key (with validation)
3. Select default model
4. Configure privacy settings
5. Test connection

### Privacy Controls

```yaml
ai:
  anonymize_data: true        # Remove student names
  consent_required: true      # Prompt before each AI call
  local_only: false          # Use local models only
  allowed_commands:          # Whitelist AI commands
    - suggest-grade
    - analyze-course
```

## Testing Strategy

### 1. Unit Tests
```go
func TestGradingService_SuggestGrade(t *testing.T) {
    mockProvider := &MockAIProvider{
        Response: `{"score": 85, "explanation": "..."}`,
    }

    service := NewGradingService(mockProvider, client, true)
    suggestion, err := service.SuggestGrade(ctx, 123, 456)

    assert.NoError(t, err)
    assert.Equal(t, 85.0, suggestion.SuggestedScore)
}
```

### 2. Integration Tests
- Test with real Canvas sandbox
- Use smaller models for cost efficiency
- Validate prompt construction

### 3. E2E Tests
```bash
# Test full command flow
canvas ai suggest-grade --assignment-id 123 --submission-id 456 --dry-run
```

## Rollout Plan

### Beta Release (v2.0-beta.1)
- Grading assistant only
- Opt-in flag: `--enable-ai-beta`
- Gather feedback from early adopters

### Stable Release (v2.0)
- All AI features
- Comprehensive documentation
- Video tutorials

### Future Enhancements (v2.1+)
- Multi-submission batch grading
- Learning outcome prediction
- Personalized student dashboards
- Integration with external AI tools

## Cost Considerations

### API Costs (Anthropic/OpenAI)
- Grade suggestion: ~$0.01-0.03 per submission
- Course analysis: ~$0.10-0.30 per course
- Quiz generation: ~$0.05-0.15 per quiz

### Cost Controls
```yaml
ai:
  monthly_budget: 50.00      # USD
  warn_at_percentage: 80
  cache_aggressively: true
```

### Local Alternative (Free)
- Ollama + Llama 3
- Slightly lower quality but zero API costs
- Good for privacy-sensitive data

## Documentation Needs

1. **User Guide**
   - `/docs/ai-features.md`
   - Examples for each command
   - Best practices

2. **Privacy Policy**
   - `/docs/ai-privacy.md`
   - Data handling transparency
   - Opt-out instructions

3. **API Documentation**
   - Provider abstraction
   - Custom provider implementation
   - Prompt engineering guide

## Success Metrics

- Adoption rate of AI commands
- Time saved (before/after surveys)
- Grading accuracy (AI vs human comparison)
- User satisfaction scores
- Feature requests and feedback

## Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| AI hallucinations | Always show confidence scores; require human review |
| Privacy concerns | Default to anonymization; local processing option |
| API costs | Budget controls; caching; local alternatives |
| Dependency on external services | Support multiple providers; local fallback |
| Bias in grading | Regular auditing; transparency in scoring |

## Timeline Summary

- **Week 1-2**: Foundation (AI abstraction, config)
- **Week 3**: Grading assistant (MVP)
- **Week 4**: Analytics insights
- **Week 5-6**: Discussion + content generation
- **Week 7**: Natural language queries
- **Week 8**: Testing, docs, polish
- **Week 9**: Beta release

## Next Steps

1. ✅ Get stakeholder buy-in on approach
2. Create feature branch: `feature/ai-integration`
3. Implement Phase 1 (foundation)
4. Build grading assistant MVP
5. User testing with real instructors
6. Iterate based on feedback

---

**Questions? Concerns? Suggestions?**

This is a living document. Please provide feedback and ideas!
