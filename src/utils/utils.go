package utils

import "fmt"

type SizeConverter struct {
   Bytes uint64
}

func (s SizeConverter) ToMB() float64 {
   return float64(s.Bytes) / (1024 * 1024)
}

func (s SizeConverter) ToGB() float64 {
   return float64(s.Bytes) / (1024 * 1024 * 1024)
}

func (s SizeConverter) ToTB() float64 {
   return float64(s.Bytes) / (1024 * 1024 * 1024 * 1024)
}

func (s SizeConverter) ToReadable() string {
   if s.Bytes < 1024 {
	   return fmt.Sprintf("%d B", s.Bytes)
   } else if s.Bytes < 1024*1024 {
	   return fmt.Sprintf("%.2f KB", float64(s.Bytes)/1024)
   } else if s.Bytes < 1024*1024*1024 {
	   return fmt.Sprintf("%.2f MB", float64(s.Bytes)/(1024*1024))
   } else if s.Bytes < 1024*1024*1024*1024 {
	   return fmt.Sprintf("%.2f GB", float64(s.Bytes)/(1024*1024*1024))
   }
   return fmt.Sprintf("%.2f TB", float64(s.Bytes)/(1024*1024*1024*1024))
}
