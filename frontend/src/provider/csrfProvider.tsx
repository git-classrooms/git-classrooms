import { authCsrfQueryOptions } from "@/api/auth";
import { apiClientOptions } from "@/lib/utils";
import { useSuspenseQuery } from "@tanstack/react-query";
import axios, { AxiosInstance } from "axios";
import React, { createContext, useMemo } from "react";

type CsrfState = {
  apiClient: AxiosInstance;
  csrfToken: string;
};

const initialState: CsrfState = {
  apiClient: axios,
  csrfToken: "",
};

const CsrfContext = createContext<CsrfState>(initialState);

export const CsrfProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { data } = useSuspenseQuery(authCsrfQueryOptions);
  const apiClient = useMemo(
    () =>
      axios.create({
        ...apiClientOptions,
        headers: {
          "X-CSRF-Token": data.csrf,
        },
      }),
    [data.csrf],
  );

  return <CsrfContext.Provider value={{ apiClient, csrfToken: data.csrf }}>{children}</CsrfContext.Provider>;
};

export const useCsrf = () => React.useContext(CsrfContext);
