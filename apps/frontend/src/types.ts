export type Template = {
  id: number;
  problem_id: number;
  language: string;
  code: string;
  metadata: any;
  created_at: string;
  updated_at: string;
};

export type Problem = {
  id: number;
  lesson_id: number;
  name: string;
  statement: string;
  difficulty: 'easy' | 'medium' | 'hard';
  time_limit: number;
  memory_limit: number;
  created_at: string;
  updated_at: string;
  template?: Template;
  languages?: string[];
};

export type Lesson = {
  id: number;
  course_id: number;
  name: string;
  description: string;
  position: number;
  content: string;
  created_at: string;
  updated_at: string;
  problems: Problem[];
};

export type Course = {
  id: number;
  name: string;
  description: string;
  tags: string[];
  status: string;
  created_at: string;
  updated_at: string;
  enrolled: boolean;
  lessons: Lesson[];
};

export type UserProfile = {
  id: number;
  login: string;
  created_at: string;
};
