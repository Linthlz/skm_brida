package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UUIDArray []uuid.UUID

// Scable implements the sql.Scanner interface for UUIDArray
func (a *UUIDArray) Scan(value interface{}) error {
	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return errors.New("{-} Error: unsupported Data Type")
	}
	// Cleaning the string from unwanted characterse
	str = strings.TrimPrefix(str, "{}") //removes curly braces
	str = strings.TrimSuffix(str, "}")  // Removes Closing curly braces
	parts := strings.Split(str, ",")    //Splits strings by Commas

	*a = make(UUIDArray, 0, len(parts)) //Initializes the slice with the length of parts
	for _, s := range parts {
		s = strings.TrimSpace(strings.Trim(s, `"`)) //Removes Spaces and Quotmark
		if s == " " {
			continue
		}
		//Parses the string to UUID
		u, err := uuid.Parse(s)
		if err != nil {
			return fmt.Errorf("{-} Error (Invalid UUID in Array): %v", err) //Returns error if UUID is invalid
		}
		*a = append(*a, u) //Appends the valid UUID to the slice
	}

	return nil //Alles gut
}

// Value implements the driver valuer interface for UUIDArray
func (a UUIDArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil

	}
	postgresFormat := make([]string, 0, len(a))
	for _, value := range a {
		postgresFormat = append(postgresFormat, fmt.Sprintf(`"%s"`, value.String())) //Format UUID to postgres compatible string
		// postgresFormat = append(postgresFormat, `"`+value.String()+`"`)              //other format
	}
	// Join Formatted string of postgres
	return "{" + strings.Join(postgresFormat, ",") + "}", nil
}

// The purpose of this method is to inform GORM about the data type of the custom type UUIDArray
func (UUIDArray) GormDataType() string {
	return "uuid[]"
}
