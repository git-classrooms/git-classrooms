"use client"

//import { zodResolver } from "@hookform/resolvers/zod";
//import { createFileRoute, useNavigate } from "@tanstack/react-router";
//import { useFieldArray, useForm } from "react-hook-form";
//import { z } from "zod";
//
//import { Button } from "@/components/ui/button"
//import {
//  Form,
//  FormControl,
//  FormDescription,
//  FormField,
//  FormItem,
//  FormLabel,
//  FormMessage,
//} from "@/components/ui/form"
//import { Input } from "@/components/ui/input"
//import { Textarea } from '@/components/ui/textarea'
//import { createFormSchema } from "@/types/classroom";
//import { createClassRoom } from "@/api/classrooms";
//import { useState } from "react";
//import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
//import { AlertCircle } from "lucide-react";
//
//export const Route = createFileRoute("/_auth/classrooms/invite")({
//  component: ClassroomsForm,
//});
//
//function ClassroomsForm() {
//  const navigate = useNavigate({ from: "/_auth/classrooms/create" })
//  const [hasError, setHasError] = useState(false)
//  const form = useForm<z.infer<typeof createFormSchema>>({
//    resolver: zodResolver(createFormSchema),
//    defaultValues: {
//      name: "",
//      description: ""
//    }
//  })
//
//  const { fields, append } = useFieldArray({
//    name: "memberEmails",
//    control: form.control,
//  })
//
//  async function onSubmit(values: z.infer<typeof createFormSchema>) {
//    try {
//      await createClassRoom(values)
//      navigate({ to: "/classrooms" })
//    } catch (error) {
//      setHasError(true)
//
//    }
//  }
//
//  return (
//    <div className="p-2">
//      <div className="flex flex-row justify-between">
//        <h1 className="text-xl font-bold">Create a classroom</h1>
//      </div>
//      <Form {...form}>
//        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
//          <FormField
//            control={form.control}
//            name="name"
//            render={({ field }) => (
//              <FormItem>
//                <FormLabel>Name</FormLabel>
//                <FormControl>
//                  <Input placeholder="Programming classroom" {...field} />
//                </FormControl>
//                <FormDescription>
//                  This is your classroom name.
//                </FormDescription>
//                <FormMessage />
//              </FormItem>
//            )}
//          />
//
//          <div>
//            {fields.map((field, index) =>
//              <FormField
//                control={form.control}
//                key={field.id}
//                name={`memberEmails.${index}`}
//                render={({ field }) => (
//                  <FormItem>
//                    <FormLabel className={index === 0 ? "" : "sr-only"}>
//                      Member emails
//                    </FormLabel>
//                    <FormControl>
//                      <Input
//                        type="email"
//                        {...field} />
//                    </FormControl>
//                    <FormMessage />
//                  </FormItem>
//                )}
//              />
//            )}
//          </div>
//          <Button onClick={() => append("")} type="button">Add Member</Button>
//
//          <FormField
//            control={form.control}
//            name="description"
//            render={({ field }) => (
//              <FormItem>
//                <FormLabel>Description</FormLabel>
//                <FormControl>
//                  <Textarea
//                    placeholder="This is my awesome ..."
//                    className="resize-none"
//                    {...field} />
//                </FormControl>
//                <FormDescription>
//                  This is the description of your classroom.
//                </FormDescription>
//                <FormMessage />
//              </FormItem>
//            )}
//          />
//          <Button type="submit">Submit</Button>
//
//          {hasError && <Alert variant="destructive">
//            <AlertCircle className="h-4 w-4" />
//            <AlertTitle>Error</AlertTitle>
//            <AlertDescription>
//              The classroom could not be created!
//            </AlertDescription>
//          </Alert>
//          }
//        </form>
//      </Form>
//    </div>
//  );
//}
//
//
