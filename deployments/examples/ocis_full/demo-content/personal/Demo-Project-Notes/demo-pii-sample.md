# Customer Intake Form — SYNTHETIC TEST DATA ONLY

> **This entire document is fabricated for testing purposes.**
> No field below refers to a real person, account, or institution.
> Names, numbers, and addresses use standard "reserved for fictional use" ranges
> (555 phone exchange, example.com/example.org email domains, test card numbers, etc.).

## Personal Information

- Full Name: Jane Q. Testperson
- Date of Birth: 1985-03-14
- Social Security Number: 000-12-3456 *(000 prefix is never issued — invalid by design)*
- Driver's License: D1234567 (State of Sampleland)
- Passport Number: X00000001 (fictional format)

## Contact Details

- Home Address: 123 Fictional Lane, Apt 4B, Testville, CA 90210
- Phone (mobile): (555) 010-2938 *(555-01xx is reserved for fiction per NANP)*
- Phone (work): (555) 010-4471
- Email: jane.testperson@example.com
- Emergency Contact: John Testperson, (555) 010-8842, john.testperson@example.org

## Financial Information

- Bank Name: First Sample National Bank
- Account Number: 0000000000
- Routing Number: 000000000
- Credit Card (Visa test number): 4111 1111 1111 1111, Exp 01/30, CVV 000
- Credit Card (Mastercard test number): 5555 5555 5555 4444, Exp 01/30, CVV 000

## Health Information (fictional)

- Patient ID: PT-000001
- Primary Physician: Dr. Alice Sampleton, Sampleland General Hospital
- Diagnosis Note: "Patient presents with fictional-condition-A, prescribed placebo 10mg daily."
- Insurance Policy Number: INS-0000-0000

## Employment

- Employer: Acme Test Corporation
- Employee ID: EMP-000123
- Salary: $000,000.00 (fictional figure)

## Notes

Use this file to exercise PII-detection features (e.g. the `ai-sensitive-data-scanner`
web extension) against a document that should trigger every major category:
name, SSN, DOB, address, phone, email, financial account/card numbers, and health data.
