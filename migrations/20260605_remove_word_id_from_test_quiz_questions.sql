-- Remove the legacy word_id relationship from quiz questions.
-- Run this in Supabase SQL Editor or your migration pipeline.

ALTER TABLE IF EXISTS test_quiz_questions
	DROP CONSTRAINT IF EXISTS fk_test_quiz_questions_word;

ALTER TABLE IF EXISTS test_quiz_questions
	DROP COLUMN IF EXISTS word_id;
