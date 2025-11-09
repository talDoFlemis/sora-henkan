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
