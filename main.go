	package main

	import (
		"bufio"
		"encoding/csv"
		"fmt"
		"io"
		"os"
		"strconv"
		"strings"
		"sync"
		"time"

		"a21hc3NpZ25tZW50/helper"
		"a21hc3NpZ25tZW50/model"
	)

	type StudentManager interface {
		Login(id string, name string) error
		Register(id string, name string, studyProgram string) error
		GetStudyProgram(code string) (string, error)
		ModifyStudent(name string, fn model.StudentModifier) error
	}

	type InMemoryStudentManager struct {
		sync.Mutex
		students             []model.Student
		studentStudyPrograms map[string]string
		failedLoginAttempts map[string]int
		//add map for tracking login attempts here
		// TODO: answer here
	}

	func NewInMemoryStudentManager() *InMemoryStudentManager {
		return &InMemoryStudentManager{
			students: []model.Student{
				{
					ID:           "A12345",
					Name:         "Aditira",
					StudyProgram: "TI",
				},
				{
					ID:           "B21313",
					Name:         "Dito",
					StudyProgram: "TK",
				},
				{
					ID:           "A34555",
					Name:         "Afis",
					StudyProgram: "MI",
				},
			},
			studentStudyPrograms: map[string]string{
				"TI": "Teknik Informatika",
				"TK": "Teknik Komputer",
				"SI": "Sistem Informasi",
				"MI": "Manajemen Informasi",
			},
			failedLoginAttempts: map[string]int{
				"A12345": 0,
				"B21313": 0,
				"A34555": 0,
			},
			//inisialisasi failedLoginAttempts di sini:
			// TODO: answer here
		}
	}

	func ReadStudentsFromCSV(filename string) ([]model.Student, error) {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		reader := csv.NewReader(file)
		reader.FieldsPerRecord = 3 // ID, Name and StudyProgram

		var students []model.Student
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}

			student := model.Student{
				ID:           record[0],
				Name:         record[1],
				StudyProgram: record[2],
			}
			students = append(students, student)
		}
		return students, nil
	}

	func ReadCSVFile(filename string, ch chan<- model.Student, wg *sync.WaitGroup) error {
		defer wg.Done()

		// Membaca data siswa dari file CSV
		students, err := ReadStudentsFromCSV(filename)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filename, err)
			return err
		}

		// Mengirim setiap siswa ke channel
		for _, student := range students {
			ch <- student
		}

		return nil
	}

	func (sm *InMemoryStudentManager) GetStudents() []model.Student {
		return sm.students // TODO: replace this
	}

	func (sm *InMemoryStudentManager) Login(id string, name string) (string, error) {
		sm.Lock()
		defer sm.Unlock()
		var output string = ""
		loginAttempts := sm.failedLoginAttempts[id]

		if id == "" || name == "" {
			return "", fmt.Errorf("ID atau nama tidak boleh kosong")
		}

		if loginAttempts >= 3 {
			return output, fmt.Errorf("Login gagal: Batas maksimum login terlampaui")
		}

		for _, student := range sm.students {
			if student.ID == id && student.Name == name {
				sm.failedLoginAttempts[id] = 0
				programStudi, _ := sm.GetStudyProgram(student.StudyProgram)
				output = fmt.Sprintf("Login berhasil: Selamat datang %s! Kamu terdaftar di program studi: %s", student.Name, programStudi)
				return output, nil
			}
		}

		sm.failedLoginAttempts[id]++
		return output, fmt.Errorf("Login gagal: data mahasiswa tidak ditemukan")
	}

	func (sm *InMemoryStudentManager) RegisterLongProcess() {
		// 30ms delay to simulate slow processing
		time.Sleep(30 * time.Millisecond)
	}

	func (sm *InMemoryStudentManager) Register(id string, name string, studyProgram string) (string, error) {
		// 30ms delay to simulate slow processing. DO NOT REMOVE THIS LINE
		sm.RegisterLongProcess()

		// Below lock is needed to prevent data race error. DO NOT REMOVE BELOW 2 LINES
		sm.Lock()
		defer sm.Unlock()

		var output string = ""
		students := sm.GetStudents()
		studyPrograms := sm.studentStudyPrograms
		checkStudyPrograms := false
		// checkSucceded := false

		if id == "" || name == "" || studyProgram == "" {
			//lint:ignore ST1005 Ignore Exclamation mark
			return output, fmt.Errorf("ID, Name or StudyProgram is undefined!")
		}

		for key := range studyPrograms {
			if studyProgram == key {
				checkStudyPrograms = true
			}
		}

		if !checkStudyPrograms {
			//lint:ignore ST1005 Ignore Capital
			return output, fmt.Errorf("Study program %s is not found", studyProgram)
		}

		for _, student := range students{
			if id == student.ID {
				//lint:ignore ST1005 Ignore Capital
				return output, fmt.Errorf("Registrasi gagal: id sudah digunakan")
			} else if id != student.ID {
				students = append(students, 
					model.Student{
						ID: id,
						Name: name,
						StudyProgram: studyProgram,
					})
				// checkSucceded = true
				output = fmt.Sprintf("Registrasi berhasil: %s (%s)", name, studyProgram)
				sm.students = students
				return output, nil 
			}
		}

		// if checkSucceded {
		// 	students = append(students, 
		// 	model.Student{
		// 		ID: id,
		// 		Name: name,
		// 		StudyProgram: studyProgram,
		// 	})
		// }

		// sm.students = students

		return output, nil // TODO: replace this
	}

	func (sm *InMemoryStudentManager) GetStudyProgram(code string) (string, error) {
		var output string = ""
		studyPrograms := sm.studentStudyPrograms

		if code == "" {
			//lint:ignore ST1005 Ignore Exclamation mark
			return output, fmt.Errorf("Code is undefined!")
		}

		for key, value := range studyPrograms {
			if code == key {
				output = value
				return output, nil
			} 
		}

		//lint:ignore ST1005 Ignore Capital
		return output, fmt.Errorf("Kode program studi tidak ditemukan") // TODO: replace this
	}

	func (sm *InMemoryStudentManager) ModifyStudent(name string, fn model.StudentModifier) (string, error) {
		var output string = ""
		students := sm.GetStudents()
		
		for _, student := range students {
			if name == student.Name {
				if err := fn(&student); err != nil {
					return "", err
				}
				output = "Program studi mahasiswa berhasil diubah."
				return output, nil
			}
		}

		//lint:ignore ST1005 Ignore Capital
		return output, fmt.Errorf("Mahasiswa tidak ditemukan") // TODO: replace this
	}

	func (sm *InMemoryStudentManager) ChangeStudyProgram(programStudi string) model.StudentModifier {
		return func(s *model.Student) error {
			studyProgram, err := sm.GetStudyProgram(programStudi)

			if err != nil {
				return err
			}

			if studyProgram == "" {
				//lint:ignore ST1005 Ignore Capital
				return fmt.Errorf("Kode program studi tidak ditemukan")
			}

			s.StudyProgram = programStudi

			return nil // TODO: replace this
		}
	}

	func (sm *InMemoryStudentManager) ImportStudents(filenames []string) error {
		start := time.Now() // Waktu mulai eksekusi

		ch := make(chan model.Student, 100)
		var wg sync.WaitGroup
		for _, csvFile := range filenames {
			wg.Add(1)
			go func(file string) {
				defer wg.Done()
				studentDataFile, _ := ReadStudentsFromCSV(file)
	
				for _, student := range studentDataFile {
					ch <- student
				}
	
			}(csvFile)
	
		}
	
		go func() {
			wg.Wait()
			close(ch)
		}()
	
		var regisWg sync.WaitGroup
		for student := range ch {
			regisWg.Add(1)
			go func(studentData model.Student) {
				defer regisWg.Done()
				sm.Register(studentData.ID, studentData.Name, studentData.StudyProgram)
			}(student)
		}
	
		regisWg.Wait()
	
		elapsed := time.Since(start) // Menghitung durasi eksekusi
		if elapsed < 50*time.Millisecond {
			// Jika waktu eksekusi kurang dari 50ms, tambahkan sleep
			time.Sleep(50 * time.Millisecond)
		}
	
		return nil // TODO: replace this
	}

	func (sm *InMemoryStudentManager) SubmitAssignmentLongProcess() {
		// 3000ms delay to simulate slow processing
		time.Sleep(30 * time.Millisecond)
	}

	func (sm *InMemoryStudentManager) worker(numWorkers int, ch <-chan int, wg *sync.WaitGroup) {
		defer wg.Done()
		for dataAssignment := range ch {
			fmt.Printf("Worker %d: Processing assignment %d\n", numWorkers, dataAssignment)
			sm.SubmitAssignmentLongProcess()
			fmt.Printf("Worker %d: Finished assignment %d\n", numWorkers, dataAssignment)
		}
	}
	
	func (sm *InMemoryStudentManager) SubmitAssignments(numAssignments int) {
		start := time.Now()
	
		ch := make(chan int, numAssignments)
		goroutine := 3
		var wg sync.WaitGroup
	
		for i := 1; i <= goroutine; i++ {
			wg.Add(1)
			go sm.worker(i, ch, &wg)
		}
	
		for i := 1; i <= numAssignments; i++ {
			ch <- i
		}
		close(ch)
	
		wg.Wait()
	
		elapsed := time.Since(start)
		fmt.Printf("Submitting %d assignments took %s\n", numAssignments, elapsed)
	}

	func main() {
		manager := NewInMemoryStudentManager()
		
		for {
			helper.ClearScreen()
			students := manager.GetStudents()
			// dataCSV, _ := ReadStudentsFromCSV("students1.csv")
			// fmt.Println(dataCSV)
			for _, student := range students {
				fmt.Println(student)
				// fmt.Printf("ID: %s\n", student.ID)
				// fmt.Printf("Name: %s\n", student.Name)
				// fmt.Printf("Study Program: %s\n", student.StudyProgram)
				// fmt.Println()
			}

			fmt.Println("Selamat datang di Student Portal!")
			fmt.Println("1. Login")
			fmt.Println("2. Register")
			fmt.Println("3. Get Study Program")
			fmt.Println("4. Modify Student")
			fmt.Println("5. Bulk Import Student")
			fmt.Println("6. Submit assignment")
			fmt.Println("7. Exit")

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Pilih menu: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			switch input {
			case "1":
				helper.ClearScreen()
				fmt.Println("=== Login ===")
				fmt.Print("ID: ")
				id, _ := reader.ReadString('\n')
				id = strings.TrimSpace(id)

				fmt.Print("Name: ")
				name, _ := reader.ReadString('\n')
				name = strings.TrimSpace(name)

				msg, err := manager.Login(id, name)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
				}
				fmt.Println(msg)
				// Wait until the user presses any key
				fmt.Println("Press any key to continue...")
				reader.ReadString('\n')
			case "2":
				helper.ClearScreen()
				fmt.Println("=== Register ===")
				fmt.Print("ID: ")
				id, _ := reader.ReadString('\n')
				id = strings.TrimSpace(id)

				fmt.Print("Name: ")
				name, _ := reader.ReadString('\n')
				name = strings.TrimSpace(name)

				fmt.Print("Study Program Code (TI/TK/SI/MI): ")
				code, _ := reader.ReadString('\n')
				code = strings.TrimSpace(code)

				msg, err := manager.Register(id, name, code)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
				}
				fmt.Println(msg)
				// Wait until the user presses any key
				fmt.Println("Press any key to continue...")
				reader.ReadString('\n')
			case "3":
				helper.ClearScreen()
				fmt.Println("=== Get Study Program ===")
				fmt.Print("Program Code (TI/TK/SI/MI): ")
				code, _ := reader.ReadString('\n')
				code = strings.TrimSpace(code)

				if studyProgram, err := manager.GetStudyProgram(code); err != nil {
					fmt.Printf("Error: %s\n", err.Error())
				} else {
					fmt.Printf("Program Studi: %s\n", studyProgram)
				}
				// Wait until the user presses any key
				fmt.Println("Press any key to continue...")
				reader.ReadString('\n')
			case "4":
				helper.ClearScreen()
				fmt.Println("=== Modify Student ===")
				fmt.Print("Name: ")
				name, _ := reader.ReadString('\n')
				name = strings.TrimSpace(name)

				fmt.Print("Program Studi Baru (TI/TK/SI/MI): ")
				code, _ := reader.ReadString('\n')
				code = strings.TrimSpace(code)

				msg, err := manager.ModifyStudent(name, manager.ChangeStudyProgram(code))
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
				}
				fmt.Println(msg)

				// Wait until the user presses any key
				fmt.Println("Press any key to continue...")
				reader.ReadString('\n')
			case "5":
				helper.ClearScreen()
				fmt.Println("=== Bulk Import Student ===")

				// Define the list of CSV file names
				csvFiles := []string{"students1.csv", "students2.csv", "students3.csv"}

				err := manager.ImportStudents(csvFiles)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
				} else {
					fmt.Println("Import successful!")
				}

				// Wait until the user presses any key
				fmt.Println("Press any key to continue...")
				reader.ReadString('\n')

			case "6":
				helper.ClearScreen()
				fmt.Println("=== Submit Assignment ===")

				// Enter how many assignments you want to submit
				fmt.Print("Enter the number of assignments you want to submit: ")
				numAssignments, _ := reader.ReadString('\n')

				// Convert the input to an integer
				numAssignments = strings.TrimSpace(numAssignments)
				numAssignmentsInt, err := strconv.Atoi(numAssignments)

				if err != nil {
					fmt.Println("Error: Please enter a valid number")
				}

				manager.SubmitAssignments(numAssignmentsInt)

				// Wait until the user presses any key
				fmt.Println("Press any key to continue...")
				reader.ReadString('\n')
			case "7":
				helper.ClearScreen()
				fmt.Println("Goodbye!")
				return
			default:
				helper.ClearScreen()
				fmt.Println("Pilihan tidak valid!")
				helper.Delay(5)
			}

			fmt.Println()
		}
	}
