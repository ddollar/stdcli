package stdcli

import (
	"context"
	"strings"
	"testing"
)

func TestArgs(t *testing.T) {
	tests := []struct {
		name      string
		validator Validator
		args      []string
		wantError bool
		errMsg    string
	}{
		{
			name:      "Args(2) with 2 args succeeds",
			validator: Args(2),
			args:      []string{"arg1", "arg2"},
			wantError: false,
		},
		{
			name:      "Args(2) with 1 arg fails",
			validator: Args(2),
			args:      []string{"arg1"},
			wantError: true,
			errMsg:    "2 args required",
		},
		{
			name:      "Args(2) with 3 args fails",
			validator: Args(2),
			args:      []string{"arg1", "arg2", "arg3"},
			wantError: true,
			errMsg:    "2 args required",
		},
		{
			name:      "Args(1) with 1 arg succeeds",
			validator: Args(1),
			args:      []string{"arg1"},
			wantError: false,
		},
		{
			name:      "Args(1) with 0 args fails with singular",
			validator: Args(1),
			args:      []string{},
			wantError: true,
			errMsg:    "1 arg required",
		},
		{
			name:      "Args(0) with 0 args succeeds",
			validator: Args(0),
			args:      []string{},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &defaultContext{
				Context: context.Background(),
				args:    tt.args,
			}

			err := tt.validator(ctx)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error to contain %q but got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestArgsMin(t *testing.T) {
	tests := []struct {
		name      string
		validator Validator
		args      []string
		wantError bool
		errMsg    string
	}{
		{
			name:      "ArgsMin(2) with 2 args succeeds",
			validator: ArgsMin(2),
			args:      []string{"arg1", "arg2"},
			wantError: false,
		},
		{
			name:      "ArgsMin(2) with 3 args succeeds",
			validator: ArgsMin(2),
			args:      []string{"arg1", "arg2", "arg3"},
			wantError: false,
		},
		{
			name:      "ArgsMin(2) with 1 arg fails",
			validator: ArgsMin(2),
			args:      []string{"arg1"},
			wantError: true,
			errMsg:    "at least 2 args required",
		},
		{
			name:      "ArgsMin(1) with 0 args fails with singular",
			validator: ArgsMin(1),
			args:      []string{},
			wantError: true,
			errMsg:    "at least 1 arg required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &defaultContext{
				Context: context.Background(),
				args:    tt.args,
			}

			err := tt.validator(ctx)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error to contain %q but got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestArgsMax(t *testing.T) {
	tests := []struct {
		name      string
		validator Validator
		args      []string
		wantError bool
		errMsg    string
	}{
		{
			name:      "ArgsMax(2) with 2 args succeeds",
			validator: ArgsMax(2),
			args:      []string{"arg1", "arg2"},
			wantError: false,
		},
		{
			name:      "ArgsMax(2) with 1 arg succeeds",
			validator: ArgsMax(2),
			args:      []string{"arg1"},
			wantError: false,
		},
		{
			name:      "ArgsMax(2) with 3 args fails",
			validator: ArgsMax(2),
			args:      []string{"arg1", "arg2", "arg3"},
			wantError: true,
			errMsg:    "no more than 2 args expected",
		},
		{
			name:      "ArgsMax(1) with 2 args fails with singular",
			validator: ArgsMax(1),
			args:      []string{"arg1", "arg2"},
			wantError: true,
			errMsg:    "no more than 1 arg expected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &defaultContext{
				Context: context.Background(),
				args:    tt.args,
			}

			err := tt.validator(ctx)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error to contain %q but got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestArgsBetween(t *testing.T) {
	tests := []struct {
		name      string
		validator Validator
		args      []string
		wantError bool
		errMsg    string
	}{
		{
			name:      "ArgsBetween(1,3) with 2 args succeeds",
			validator: ArgsBetween(1, 3),
			args:      []string{"arg1", "arg2"},
			wantError: false,
		},
		{
			name:      "ArgsBetween(1,3) with 1 arg succeeds",
			validator: ArgsBetween(1, 3),
			args:      []string{"arg1"},
			wantError: false,
		},
		{
			name:      "ArgsBetween(1,3) with 3 args succeeds",
			validator: ArgsBetween(1, 3),
			args:      []string{"arg1", "arg2", "arg3"},
			wantError: false,
		},
		{
			name:      "ArgsBetween(1,3) with 0 args fails",
			validator: ArgsBetween(1, 3),
			args:      []string{},
			wantError: true,
			errMsg:    "at least 1 arg required",
		},
		{
			name:      "ArgsBetween(1,3) with 4 args fails",
			validator: ArgsBetween(1, 3),
			args:      []string{"arg1", "arg2", "arg3", "arg4"},
			wantError: true,
			errMsg:    "no more than 3 args expected",
		},
		{
			name:      "ArgsBetween(0,1) with 0 args succeeds",
			validator: ArgsBetween(0, 1),
			args:      []string{},
			wantError: false,
		},
		{
			name:      "ArgsBetween(0,1) with 1 arg succeeds",
			validator: ArgsBetween(0, 1),
			args:      []string{"arg1"},
			wantError: false,
		},
		{
			name:      "ArgsBetween(0,1) with 2 args fails",
			validator: ArgsBetween(0, 1),
			args:      []string{"arg1", "arg2"},
			wantError: true,
			errMsg:    "no more than 1 arg expected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &defaultContext{
				Context: context.Background(),
				args:    tt.args,
			}

			err := tt.validator(ctx)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error to contain %q but got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidatorsInCommand(t *testing.T) {
	tests := []struct {
		name      string
		validator Validator
		args      []string
		wantError bool
		errMsg    string
	}{
		{
			name:      "command with Args(2) validator accepts 2 args",
			validator: Args(2),
			args:      []string{"arg1", "arg2"},
			wantError: false,
		},
		{
			name:      "command with Args(2) validator rejects 1 arg",
			validator: Args(2),
			args:      []string{"arg1"},
			wantError: true,
			errMsg:    "2 args required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New("testapp", "1.0.0")

			handlerCalled := false
			e.Command("test", "test command", func(ctx Context) error {
				handlerCalled = true
				return nil
			}, CommandOptions{
				Validate: tt.validator,
			})

			// Find the test command
			var cmd *Command
			for i, c := range e.Commands {
				if len(c.Command) > 0 && c.Command[0] == "test" {
					cmd = &e.Commands[i]
					break
				}
			}

			if cmd == nil {
				t.Fatal("test command not found")
			}

			err := cmd.ExecuteContext(context.Background(), tt.args)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error to contain %q but got %q", tt.errMsg, err.Error())
				}
				if handlerCalled {
					t.Errorf("handler should not be called when validation fails")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !handlerCalled {
					t.Errorf("handler should be called when validation succeeds")
				}
			}
		})
	}
}
