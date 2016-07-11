package help

import (
  "os"
  //"io"
  "io/ioutil"
  "fmt"
  "testing"

  . "github.com/smartystreets/goconvey/convey"
)

var helpText = `
base =  <<EOF
This is the base text
EOF

test1 = <<EOF
This is test1
EOF
test2 = <<EOF
This is test2EOF

test1.test2 = <<EOF
This is test1.test2
EOF
test1.test2.test3 = "This is test1.test2.test3"

`

var testData = []struct {
  Key string
  Value string
} {
  {
    Base,
    "This is the base text\n",
  },
  {
    "test1",
    "This is test1\n",
  },
  {
    "test2",
    "This is test2",
  },
  {
    "test1.test2",
    "This is test1.test2\n",
  },
  {
    "test1.test2.test3",
    "This is test1.test2.test3",
  },

}

// guarantee that a tmp directory is uses.  Precedence:
// GoLang TempDir()
// ENV TMPDIR

func GetTmpDir() string {
  tmpDir := os.TempDir()
  if tmpDir == "" {
    tmpDir = os.Getenv("TMPDIR")
  }
  if tmpDir == "" {
    tmpDir = "/tmp"
  }
  if tmpDir == "" {
    tmpDir = "."
  }
  fi, err := os.Stat(tmpDir)
  if err != nil {
      panic("Could not stat temp folder!")
  }

  if !fi.IsDir() {
    panic("Could not find a temp folder!")
  }
  // we should test the folder for write permission
  // permissions := fi.Perm()

  return tmpDir
}

func CreateHelpFile(txt string) string {
  tmpFile, err := ioutil.TempFile(GetTmpDir(), "gorest_test_")
  if err != nil {
    // Yikes!  Could not write to file!
    panic("Could not create tmpFile!")
  }
  _, err = tmpFile.WriteString(txt)
  if err != nil {
    // Yikes!  Could not write to file!
    panic("Could not write to tmpFile!")
  }
  tmpFile.Close()
  return tmpFile.Name()
}

func RemoveHelpFile(path string) {
  err := os.Remove(path)
  if err != nil {
    panic(err)
  }
}

func TestParseHelpText(t *testing.T) {
  Convey("Given the help text", t, func() {
    // Create tmp file
    out, err := ParseHelpText(helpText)
    if err != nil {
      panic(err)
    }
    for _, d := range testData {
      Convey(fmt.Sprintf("Given the key [%s]", d.Key), func() {
        txt := out.Get(d.Key)
        Convey(fmt.Sprintf("The value should be [%s]", d.Value), func() {
          So(txt, ShouldEqual, d.Value)
        })
      })
    }
  })
}

func TestParseHelpFile(t *testing.T) {

  Convey("Given a help file", t, func() {
    // Create tmp file
    path := CreateHelpFile(helpText)
    out, err := ParseHelpFile(path)
    if err != nil {
      panic(err)
    }
    for _, d := range testData {
      Convey(fmt.Sprintf("Given the key [%s]", d.Key), func() {
        txt := out.Get(d.Key)
        Convey(fmt.Sprintf("The value should be [%s]", d.Value), func() {
          So(txt, ShouldEqual, d.Value)
        })
      })
    }
    Reset(func() {
      // Get rid of tmp file
      RemoveHelpFile(path)
    })
  })
}
