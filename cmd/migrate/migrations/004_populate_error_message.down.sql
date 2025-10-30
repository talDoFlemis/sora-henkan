-- Revert error_message population (set back to NULL)
UPDATE images SET error_message = NULL WHERE error_message = '';
