# ============================================================================
# IAM Configuration for AWS Academy / Learner Lab Environments
# ============================================================================
# 
# In AWS Academy and Learner Lab environments, IAM permissions are restricted:
# - You CANNOT create new IAM roles or policies
# - A pre-existing role named "LabRole" is provided
# - A pre-existing instance profile named "LabInstanceProfile" is provided
# 
# The LabRole has broad permissions to access AWS services including:
# - S3 (read/write objects)
# - SQS (send/receive/delete messages)
# - CloudWatch (logs and metrics)
# - RDS (database connections)
# - And many other AWS services
#
# This configuration uses data sources to reference these existing resources
# instead of trying to create new ones.
# ============================================================================

# Data source for existing LabRole
# This role is pre-created in AWS Academy/Learner Lab environments
# and has permissions to access S3, SQS, CloudWatch, and other AWS services
data "aws_iam_role" "lab_role" {
  name = "LabRole"
}

# Data source for existing LabInstanceProfile
# This instance profile is pre-created and associated with LabRole
data "aws_iam_instance_profile" "lab_profile" {
  name = "LabInstanceProfile"
}
