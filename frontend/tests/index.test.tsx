import { expect, test, describe } from "vitest";
import { createAssignmentFormSchema } from "../src/types/assignments";

describe("test create Assignmet Form Schema", () => {
  test("errors on null value ", () => {
    expect(() => createAssignmentFormSchema.parse(undefined)).toThrowError();
    expect(() => createAssignmentFormSchema.parse(null)).toThrowError();
  });

  test("errors on empty object", () => {
    expect(() => createAssignmentFormSchema.parse({})).toThrowError();
  });

  const validForm = {
    name: "Test",
    description: "Test",
    templateProjectId: 1,
    dueDate: new Date("2024-05-30T19:37:40.787Z"),
  } as const;

  test("errors on to small name", () => {
    expect(() => createAssignmentFormSchema.parse({ ...validForm, name: undefined })).toThrowError();
    expect(() => createAssignmentFormSchema.parse({ ...validForm, name: null })).toThrowError();
    expect(() => createAssignmentFormSchema.parse({ ...validForm, name: "" })).toThrowError();
    expect(() => createAssignmentFormSchema.parse({ ...validForm, name: "12" })).toThrowError();
  });

  test("errors on to invalid name", () => {
    expect(() => createAssignmentFormSchema.parse({ ...validForm, name: 3 })).toThrowError();
    expect(() => createAssignmentFormSchema.parse({ ...validForm, name: "Hallo!" })).toThrowError();
    expect(() => createAssignmentFormSchema.parse({ ...validForm, name: "Hallo⽬" })).toThrowError();
  });

  test("errors on to small description", () => {
    expect(() => createAssignmentFormSchema.parse({ description: undefined })).toThrowError();
    expect(() => createAssignmentFormSchema.parse({ description: null })).toThrowError();
    expect(() => createAssignmentFormSchema.parse({ description: 3 })).toThrowError();
    expect(() => createAssignmentFormSchema.parse({ description: "" })).toThrowError();
    expect(() => createAssignmentFormSchema.parse({ description: "12" })).toThrowError();
  });

  test("errors on invalid Date", () => {
    expect(() => createAssignmentFormSchema.parse({ ...validForm, dueDate: undefined })).toThrowError();
    expect(createAssignmentFormSchema.parse({ ...validForm, dueDate: null })).toEqual({
      ...validForm,
      dueDate: new Date("1970-01-01T00:00:00.000Z"),
    });
    expect(() => createAssignmentFormSchema.parse({ ...validForm, dueDate: "Hallo" })).toThrowError();
  });

  test("valid Form", () => {
    expect(createAssignmentFormSchema.parse(validForm)).toEqual(validForm);
    expect(createAssignmentFormSchema.parse({ ...validForm, dueDate: "2024-05-30T19:37:40.787Z" })).toEqual(validForm);
  });

  test("valid Form with special characters in Name", () => {
    const valid = {
      ...validForm,
      name: "Hallo❤ wie geht es dir",
    } as const;
    expect(createAssignmentFormSchema.parse(valid)).toEqual(valid);

    const valid2 = {
      ...validForm,
      name: "Hallo +-_ äüö wie geht es dir",
    } as const;
    expect(createAssignmentFormSchema.parse(valid2)).toEqual(valid2);
  });
});
