import { authCsrfQueryOptions } from "@/api/auth";
import { useSuspenseQuery } from "@tanstack/react-query";
import React, { createContext } from "react";

type CsrfState = {
  csrfToken: string;
};

const initialState: CsrfState = {
  csrfToken: "",
};

const CsrfContext = createContext<CsrfState>(initialState);

export const CsrfProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { data } = useSuspenseQuery(authCsrfQueryOptions);

  return <CsrfContext.Provider value={{ csrfToken: data.csrf }}>{children}</CsrfContext.Provider>;
};

export const useCsrf = () => React.useContext(CsrfContext);
