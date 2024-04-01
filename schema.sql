-- Create Teachers Table
CREATE TABLE IF NOT EXISTS teachers (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL
);

-- Create Students Table
CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    is_suspended BOOLEAN NOT NULL DEFAULT false
);

-- Create TeacherStudents Junction Table
CREATE TABLE IF NOT EXISTS teacher_students (
    teacher_id INT NOT NULL,
    student_id INT NOT NULL,
   	PRIMARY KEY (teacher_id, student_id),
    CONSTRAINT fk_teacher
        FOREIGN KEY(teacher_id) 
        REFERENCES teachers(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_student
        FOREIGN KEY(student_id) 
        REFERENCES students(id)
        ON DELETE CASCADE
);
