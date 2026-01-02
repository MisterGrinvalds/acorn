#!/usr/bin/env bash
# AWS Cloud Automation Module

# Module configuration
readonly AWS_TEMPLATES_DIR="$AUTO_HOME/cloud/aws/templates"
readonly AWS_CONFIGS_DIR="$AUTO_HOME/cloud/aws/configs"
readonly AWS_SCRIPTS_DIR="$AUTO_HOME/cloud/aws/scripts"

# Initialize directories
mkdir -p "$AWS_TEMPLATES_DIR" "$AWS_CONFIGS_DIR" "$AWS_SCRIPTS_DIR"

# Help function for AWS module
aws_help() {
    cat << EOF
AWS Cloud Automation

USAGE:
    auto aws <command> [options]

COMMANDS:
    auth                          Authentication and profile management
    ec2 <action>                  EC2 instance management
    s3 <action>                   S3 bucket operations
    lambda <action>               Lambda function management
    ecs <action>                  ECS cluster and service management
    rds <action>                  RDS database operations
    iam <action>                  IAM user and role management
    cloudformation <action>       CloudFormation stack management
    logs <action>                 CloudWatch logs operations
    cost                          Cost analysis and optimization
    security                      Security audit and compliance
    backup                        Backup and disaster recovery

EXAMPLES:
    auto aws auth setup
    auto aws ec2 list --region us-east-1
    auto aws s3 create my-bucket
    auto aws lambda deploy my-function
    auto aws ecs deploy my-app production
    auto aws rds create my-db mysql
    auto aws cloudformation deploy my-stack
    auto aws cost analyze --service ec2
EOF
}

# Utility functions
require_aws_cli() {
    require_command aws
}

get_aws_profile() {
    aws configure list-profiles 2>/dev/null | head -1 || echo "default"
}

get_aws_region() {
    aws configure get region 2>/dev/null || echo "us-east-1"
}

# Authentication and profile management
aws_auth_setup() {
    require_aws_cli
    
    log "INFO" "Setting up AWS authentication..."
    
    # Check if AWS CLI is configured
    if ! aws sts get-caller-identity >/dev/null 2>&1; then
        log "INFO" "AWS CLI not configured. Starting configuration..."
        aws configure
    else
        log "INFO" "AWS CLI already configured"
        aws sts get-caller-identity --output table
    fi
    
    # Setup additional profiles if needed
    if confirm "Setup additional AWS profiles?"; then
        aws_auth_add_profile
    fi
    
    log "SUCCESS" "AWS authentication setup complete"
}

aws_auth_add_profile() {
    echo -n "Enter profile name: "
    read -r profile_name
    
    echo -n "Enter AWS Access Key ID: "
    read -r access_key
    
    echo -n "Enter AWS Secret Access Key: "
    read -rs secret_key
    echo ""
    
    echo -n "Enter default region [us-east-1]: "
    read -r region
    region=${region:-us-east-1}
    
    aws configure set aws_access_key_id "$access_key" --profile "$profile_name"
    aws configure set aws_secret_access_key "$secret_key" --profile "$profile_name"
    aws configure set region "$region" --profile "$profile_name"
    
    log "SUCCESS" "Profile '$profile_name' created"
}

aws_auth_list_profiles() {
    require_aws_cli
    
    log "INFO" "AWS Profiles:"
    aws configure list-profiles | while read -r profile; do
        local region=$(aws configure get region --profile "$profile" 2>/dev/null || echo "not-set")
        printf "  %-15s Region: %s\n" "$profile" "$region"
    done
}

# EC2 operations
aws_ec2_list() {
    require_aws_cli
    local region="${1:-$(get_aws_region)}"
    
    log "INFO" "Listing EC2 instances in region: $region"
    
    aws ec2 describe-instances \
        --region "$region" \
        --query 'Reservations[*].Instances[*].[InstanceId,InstanceType,State.Name,PublicIpAddress,Tags[?Key==`Name`].Value|[0]]' \
        --output table
}

aws_ec2_create() {
    require_aws_cli
    local instance_type="${1:-t3.micro}"
    local ami_id="${2:-ami-0c02fb55956c7d316}" # Amazon Linux 2
    local key_name="$3"
    local security_group="${4:-default}"
    
    if [ -z "$key_name" ]; then
        log "ERROR" "Key pair name is required"
        exit 1
    fi
    
    log "INFO" "Creating EC2 instance..."
    
    local instance_id=$(aws ec2 run-instances \
        --image-id "$ami_id" \
        --instance-type "$instance_type" \
        --key-name "$key_name" \
        --security-groups "$security_group" \
        --query 'Instances[0].InstanceId' \
        --output text)
    
    if [ "$instance_id" != "None" ]; then
        log "SUCCESS" "EC2 instance created: $instance_id"
        
        # Wait for instance to be running
        log "INFO" "Waiting for instance to be running..."
        aws ec2 wait instance-running --instance-ids "$instance_id"
        
        # Get public IP
        local public_ip=$(aws ec2 describe-instances \
            --instance-ids "$instance_id" \
            --query 'Reservations[0].Instances[0].PublicIpAddress' \
            --output text)
        
        log "SUCCESS" "Instance is running. Public IP: $public_ip"
    else
        log "ERROR" "Failed to create EC2 instance"
        exit 1
    fi
}

aws_ec2_terminate() {
    require_aws_cli
    local instance_id="$1"
    
    if [ -z "$instance_id" ]; then
        log "ERROR" "Instance ID is required"
        exit 1
    fi
    
    if confirm "Terminate instance $instance_id?"; then
        aws ec2 terminate-instances --instance-ids "$instance_id"
        log "SUCCESS" "Instance $instance_id terminated"
    fi
}

# S3 operations
aws_s3_list() {
    require_aws_cli
    
    log "INFO" "Listing S3 buckets:"
    aws s3 ls
}

aws_s3_create() {
    require_aws_cli
    local bucket_name="$1"
    local region="${2:-$(get_aws_region)}"
    
    if [ -z "$bucket_name" ]; then
        log "ERROR" "Bucket name is required"
        exit 1
    fi
    
    log "INFO" "Creating S3 bucket: $bucket_name"
    
    if [ "$region" = "us-east-1" ]; then
        aws s3 mb "s3://$bucket_name"
    else
        aws s3 mb "s3://$bucket_name" --region "$region"
    fi
    
    # Enable versioning
    aws s3api put-bucket-versioning \
        --bucket "$bucket_name" \
        --versioning-configuration Status=Enabled
    
    # Enable encryption
    aws s3api put-bucket-encryption \
        --bucket "$bucket_name" \
        --server-side-encryption-configuration '{
            "Rules": [{
                "ApplyServerSideEncryptionByDefault": {
                    "SSEAlgorithm": "AES256"
                }
            }]
        }'
    
    log "SUCCESS" "S3 bucket created with versioning and encryption enabled"
}

aws_s3_sync() {
    require_aws_cli
    local source="$1"
    local destination="$2"
    local delete_flag="${3:-false}"
    
    if [ -z "$source" ] || [ -z "$destination" ]; then
        log "ERROR" "Source and destination are required"
        exit 1
    fi
    
    local sync_args=("--region" "$(get_aws_region)")
    [ "$delete_flag" = "true" ] && sync_args+=("--delete")
    
    log "INFO" "Syncing $source to $destination"
    aws s3 sync "$source" "$destination" "${sync_args[@]}"
    log "SUCCESS" "Sync completed"
}

# Lambda operations
aws_lambda_list() {
    require_aws_cli
    
    log "INFO" "Listing Lambda functions:"
    aws lambda list-functions \
        --query 'Functions[*].[FunctionName,Runtime,LastModified,CodeSize]' \
        --output table
}

aws_lambda_deploy() {
    require_aws_cli
    local function_name="$1"
    local zip_file="$2"
    local handler="${3:-index.handler}"
    local runtime="${4:-python3.9}"
    local role_arn="$5"
    
    if [ -z "$function_name" ] || [ -z "$zip_file" ]; then
        log "ERROR" "Function name and zip file are required"
        exit 1
    fi
    
    if [ ! -f "$zip_file" ]; then
        log "ERROR" "Zip file not found: $zip_file"
        exit 1
    fi
    
    # Check if function exists
    if aws lambda get-function --function-name "$function_name" >/dev/null 2>&1; then
        log "INFO" "Updating existing Lambda function: $function_name"
        aws lambda update-function-code \
            --function-name "$function_name" \
            --zip-file "fileb://$zip_file"
    else
        log "INFO" "Creating new Lambda function: $function_name"
        
        if [ -z "$role_arn" ]; then
            log "ERROR" "IAM role ARN is required for new functions"
            exit 1
        fi
        
        aws lambda create-function \
            --function-name "$function_name" \
            --runtime "$runtime" \
            --role "$role_arn" \
            --handler "$handler" \
            --zip-file "fileb://$zip_file"
    fi
    
    log "SUCCESS" "Lambda function deployed: $function_name"
}

aws_lambda_invoke() {
    require_aws_cli
    local function_name="$1"
    local payload="${2:-{}}"
    local output_file="/tmp/lambda-output.json"
    
    if [ -z "$function_name" ]; then
        log "ERROR" "Function name is required"
        exit 1
    fi
    
    log "INFO" "Invoking Lambda function: $function_name"
    
    aws lambda invoke \
        --function-name "$function_name" \
        --payload "$payload" \
        "$output_file"
    
    log "INFO" "Function output:"
    cat "$output_file"
}

# ECS operations
aws_ecs_list_clusters() {
    require_aws_cli
    
    log "INFO" "Listing ECS clusters:"
    aws ecs list-clusters --query 'clusterArns[*]' --output table
}

aws_ecs_create_cluster() {
    require_aws_cli
    local cluster_name="$1"
    
    if [ -z "$cluster_name" ]; then
        log "ERROR" "Cluster name is required"
        exit 1
    fi
    
    log "INFO" "Creating ECS cluster: $cluster_name"
    aws ecs create-cluster --cluster-name "$cluster_name"
    log "SUCCESS" "ECS cluster created: $cluster_name"
}

aws_ecs_deploy_service() {
    require_aws_cli
    local service_name="$1"
    local task_definition="$2"
    local cluster="${3:-default}"
    local desired_count="${4:-1}"
    
    if [ -z "$service_name" ] || [ -z "$task_definition" ]; then
        log "ERROR" "Service name and task definition are required"
        exit 1
    fi
    
    log "INFO" "Deploying ECS service: $service_name"
    
    # Check if service exists
    if aws ecs describe-services --cluster "$cluster" --services "$service_name" --query 'services[0].serviceName' --output text 2>/dev/null | grep -q "$service_name"; then
        log "INFO" "Updating existing service"
        aws ecs update-service \
            --cluster "$cluster" \
            --service "$service_name" \
            --task-definition "$task_definition" \
            --desired-count "$desired_count"
    else
        log "INFO" "Creating new service"
        aws ecs create-service \
            --cluster "$cluster" \
            --service-name "$service_name" \
            --task-definition "$task_definition" \
            --desired-count "$desired_count"
    fi
    
    log "SUCCESS" "ECS service deployed: $service_name"
}

# RDS operations
aws_rds_list() {
    require_aws_cli
    
    log "INFO" "Listing RDS instances:"
    aws rds describe-db-instances \
        --query 'DBInstances[*].[DBInstanceIdentifier,DBInstanceClass,Engine,DBInstanceStatus,Endpoint.Address]' \
        --output table
}

aws_rds_create() {
    require_aws_cli
    local db_identifier="$1"
    local engine="${2:-mysql}"
    local instance_class="${3:-db.t3.micro}"
    local allocated_storage="${4:-20}"
    local master_username="${5:-admin}"
    local master_password="$6"
    
    if [ -z "$db_identifier" ]; then
        log "ERROR" "DB identifier is required"
        exit 1
    fi
    
    if [ -z "$master_password" ]; then
        echo -n "Enter master password: "
        read -rs master_password
        echo ""
    fi
    
    log "INFO" "Creating RDS instance: $db_identifier"
    
    aws rds create-db-instance \
        --db-instance-identifier "$db_identifier" \
        --db-instance-class "$instance_class" \
        --engine "$engine" \
        --allocated-storage "$allocated_storage" \
        --master-username "$master_username" \
        --master-user-password "$master_password" \
        --no-publicly-accessible \
        --storage-encrypted
    
    log "SUCCESS" "RDS instance creation initiated: $db_identifier"
    log "INFO" "Instance will take several minutes to become available"
}

# CloudFormation operations
aws_cloudformation_deploy() {
    require_aws_cli
    local stack_name="$1"
    local template_file="$2"
    local parameters_file="$3"
    
    if [ -z "$stack_name" ] || [ -z "$template_file" ]; then
        log "ERROR" "Stack name and template file are required"
        exit 1
    fi
    
    if [ ! -f "$template_file" ]; then
        log "ERROR" "Template file not found: $template_file"
        exit 1
    fi
    
    local deploy_args=(
        "--stack-name" "$stack_name"
        "--template-body" "file://$template_file"
        "--capabilities" "CAPABILITY_IAM"
    )
    
    if [ -n "$parameters_file" ] && [ -f "$parameters_file" ]; then
        deploy_args+=("--parameter-overrides" "file://$parameters_file")
    fi
    
    log "INFO" "Deploying CloudFormation stack: $stack_name"
    
    aws cloudformation deploy "${deploy_args[@]}"
    log "SUCCESS" "CloudFormation stack deployed: $stack_name"
}

aws_cloudformation_list() {
    require_aws_cli
    
    log "INFO" "Listing CloudFormation stacks:"
    aws cloudformation list-stacks \
        --stack-status-filter CREATE_COMPLETE UPDATE_COMPLETE \
        --query 'StackSummaries[*].[StackName,StackStatus,CreationTime]' \
        --output table
}

# Cost analysis
aws_cost_analyze() {
    require_aws_cli
    local service="${1:-}"
    local time_period="${2:-MONTHLY}"
    
    log "INFO" "Analyzing AWS costs..."
    
    local start_date=$(date -d "last month" +%Y-%m-01)
    local end_date=$(date +%Y-%m-01)
    
    local query='{
        "TimePeriod": {
            "Start": "'$start_date'",
            "End": "'$end_date'"
        },
        "Granularity": "'$time_period'",
        "Metrics": ["BlendedCost"],
        "GroupBy": [
            {
                "Type": "DIMENSION",
                "Key": "SERVICE"
            }
        ]
    }'
    
    if [ -n "$service" ]; then
        query=$(echo "$query" | jq --arg service "$service" '.GroupBy[0].Key = "USAGE_TYPE" | .Filters = {"Dimensions": {"Key": "SERVICE", "Values": [$service]}}')
    fi
    
    aws ce get-cost-and-usage --cli-input-json "$query" \
        --query 'ResultsByTime[*].Groups[*].[Keys[0],Metrics.BlendedCost.Amount]' \
        --output table
}

# Main AWS module function
aws_main() {
    local command="${1:-help}"
    shift || true
    
    case "$command" in
        "help"|"-h"|"--help")
            aws_help
            ;;
        "auth")
            local action="${1:-setup}"
            case "$action" in
                "setup") aws_auth_setup ;;
                "add-profile") aws_auth_add_profile ;;
                "list-profiles") aws_auth_list_profiles ;;
                *) log "ERROR" "Unknown auth action: $action" ;;
            esac
            ;;
        "ec2")
            local action="$1"
            case "$action" in
                "list") aws_ec2_list "$2" ;;
                "create") aws_ec2_create "$2" "$3" "$4" "$5" ;;
                "terminate") aws_ec2_terminate "$2" ;;
                *) log "ERROR" "Unknown ec2 action: $action" ;;
            esac
            ;;
        "s3")
            local action="$1"
            case "$action" in
                "list") aws_s3_list ;;
                "create") aws_s3_create "$2" "$3" ;;
                "sync") aws_s3_sync "$2" "$3" "$4" ;;
                *) log "ERROR" "Unknown s3 action: $action" ;;
            esac
            ;;
        "lambda")
            local action="$1"
            case "$action" in
                "list") aws_lambda_list ;;
                "deploy") aws_lambda_deploy "$2" "$3" "$4" "$5" "$6" ;;
                "invoke") aws_lambda_invoke "$2" "$3" ;;
                *) log "ERROR" "Unknown lambda action: $action" ;;
            esac
            ;;
        "ecs")
            local action="$1"
            case "$action" in
                "list-clusters") aws_ecs_list_clusters ;;
                "create-cluster") aws_ecs_create_cluster "$2" ;;
                "deploy") aws_ecs_deploy_service "$2" "$3" "$4" "$5" ;;
                *) log "ERROR" "Unknown ecs action: $action" ;;
            esac
            ;;
        "rds")
            local action="$1"
            case "$action" in
                "list") aws_rds_list ;;
                "create") aws_rds_create "$2" "$3" "$4" "$5" "$6" "$7" ;;
                *) log "ERROR" "Unknown rds action: $action" ;;
            esac
            ;;
        "cloudformation"|"cf")
            local action="$1"
            case "$action" in
                "deploy") aws_cloudformation_deploy "$2" "$3" "$4" ;;
                "list") aws_cloudformation_list ;;
                *) log "ERROR" "Unknown cloudformation action: $action" ;;
            esac
            ;;
        "cost")
            aws_cost_analyze "$1" "$2"
            ;;
        *)
            log "ERROR" "Unknown command: $command"
            aws_help
            exit 1
            ;;
    esac
}