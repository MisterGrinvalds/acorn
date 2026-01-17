# Component Management
# List, status, create, and validate components

.PHONY: component-list component-status component-new component-validate
.PHONY: test-components test-component-loader test-component-deps

# List components
component-list: ## List all available components
	@echo "Available Components"
	@echo "===================="
	@echo ""
	@if [ -d "$(COMPONENTS_DIR)" ]; then \
		for dir in $(COMPONENTS_DIR)/*/; do \
			name=$$(basename "$$dir"); \
			if [ -f "$$dir/component.yaml" ]; then \
				desc=$$(yq -r '.description // "No description"' "$$dir/component.yaml" 2>/dev/null || echo "No description"); \
				printf "  %-15s %s\n" "$$name" "$$desc"; \
			else \
				printf "  %-15s (no component.yaml)\n" "$$name"; \
			fi; \
		done; \
	else \
		echo "No components directory found"; \
	fi

# Component status
component-status: ## Show component health and loading status
	@echo "Component Health Status"
	@echo "======================="
	@echo ""
	@if [ -d "$(COMPONENTS_DIR)" ]; then \
		for dir in $(COMPONENTS_DIR)/*/; do \
			name=$$(basename "$$dir"); \
			echo "$$name:"; \
			[ -f "$$dir/component.yaml" ] && echo "  component.yaml: ok" || echo "  component.yaml: missing"; \
			[ -f "$$dir/env.sh" ] && echo "  env.sh: ok" || echo "  env.sh: missing"; \
			[ -f "$$dir/aliases.sh" ] && echo "  aliases.sh: ok" || true; \
			[ -f "$$dir/functions.sh" ] && echo "  functions.sh: ok" || true; \
		done; \
	fi

# Create new component
component-new: ## Create new component (usage: make component-new NAME=mycomponent)
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make component-new NAME=mycomponent"; exit 1; \
	fi
	@if [ -d "$(COMPONENTS_DIR)/$(NAME)" ]; then \
		echo "Component '$(NAME)' already exists"; exit 1; \
	fi
	@mkdir -p "$(COMPONENTS_DIR)/$(NAME)"
	@echo "name: $(NAME)" > "$(COMPONENTS_DIR)/$(NAME)/component.yaml"
	@echo "description: $(NAME) component" >> "$(COMPONENTS_DIR)/$(NAME)/component.yaml"
	@echo "version: 1.0.0" >> "$(COMPONENTS_DIR)/$(NAME)/component.yaml"
	@echo "" >> "$(COMPONENTS_DIR)/$(NAME)/component.yaml"
	@echo "# env:" >> "$(COMPONENTS_DIR)/$(NAME)/component.yaml"
	@echo "#   MY_VAR: value" >> "$(COMPONENTS_DIR)/$(NAME)/component.yaml"
	@echo "" >> "$(COMPONENTS_DIR)/$(NAME)/component.yaml"
	@echo "# aliases:" >> "$(COMPONENTS_DIR)/$(NAME)/component.yaml"
	@echo "#   myalias: 'my command'" >> "$(COMPONENTS_DIR)/$(NAME)/component.yaml"
	@echo "Created component: $(COMPONENTS_DIR)/$(NAME)"

# Validate components
component-validate: ## Validate all component.yaml files
	@echo "Validating Components"
	@echo "====================="
	@if [ -d "$(COMPONENTS_DIR)" ]; then \
		for dir in $(COMPONENTS_DIR)/*/; do \
			name=$$(basename "$$dir"); \
			if [ -f "$$dir/component.yaml" ]; then \
				if yq '.' "$$dir/component.yaml" >/dev/null 2>&1; then \
					echo "  $$name: valid"; \
				else \
					echo "  $$name: invalid YAML"; \
				fi; \
			else \
				echo "  $$name: missing component.yaml"; \
			fi; \
		done; \
	fi

# Component tests
test-components: test-component-loader test-component-deps ## Test component system

test-component-loader: ## Test component loader functionality
	@echo "Testing component loader..."
	@bash -c '\
		export DOTFILES_ROOT="$(DOTFILES_DIR)" IS_INTERACTIVE=true; \
		source core/bootstrap.sh; \
		declare -f load_component >/dev/null && echo "load_component: defined" || echo "load_component: not defined"; \
	'
	@if [ -d "$(COMPONENTS_DIR)" ]; then \
		for dir in $(COMPONENTS_DIR)/*/; do \
			name=$$(basename "$$dir"); \
			if [ -f "$$dir/env.sh" ]; then \
				if bash -n "$$dir/env.sh" 2>/dev/null; then \
					echo "  $$name/env.sh: ok"; \
				else \
					echo "  $$name/env.sh: syntax error"; \
				fi; \
			fi; \
		done; \
	fi

test-component-deps: ## Test component dependency resolution
	@echo "Checking component dependencies..."
	@if [ -d "$(COMPONENTS_DIR)" ]; then \
		for dir in $(COMPONENTS_DIR)/*/; do \
			name=$$(basename "$$dir"); \
			if [ -f "$$dir/component.yaml" ]; then \
				deps=$$(yq -r '.requires.components // [] | .[]' "$$dir/component.yaml" 2>/dev/null); \
				for dep in $$deps; do \
					if [ ! -d "$(COMPONENTS_DIR)/$$dep" ]; then \
						echo "  $$name requires missing: $$dep"; \
					fi; \
				done; \
			fi; \
		done; \
	fi
