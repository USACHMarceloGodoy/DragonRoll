<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="w-full max-w-md bg-white p-6 rounded-2xl shadow-md">
      <h2 class="text-2xl font-bold text-center mb-6">Iniciar sesión</h2>

      <form class="space-y-4" @submit.prevent="onSubmit">
        <div>
          <label class="block text-sm font-medium text-gray-700"
            >Nombre de Usuario</label
          >
          <InputText v-model="username" type="username" class="w-full mt-1" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700"
            >Contraseña</label
          >
          <Password
            v-model="password"
            toggle-mask
            :feedback="false"
            class="w-full mt-1"
          />
        </div>
        <Button
          label="Ingresar"
          type="submit"
          class="w-full"
          severity="primary"
        />
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import InputText from "primevue/inputtext";
import Password from "primevue/password";
import Button from "primevue/button";

definePageMeta({
  title: "Iniciar sesión",
  middleware: ["guest-only"],
});

const username = ref("");
const password = ref("");
const { login } = useAuth();
const { showSuccess, showError } = useAddToast();

const onSubmit = async () => {
  const credentials = {
    user: username.value,
    password: password.value,
  };
  login(credentials)
    .then(() => {
      // Redirigir a la página principal o donde sea necesario
      useRouter().push("/");
      showSuccess("Inicio de sesión exitoso");
    })
    .catch((error) => {
      console.error(error);
      showError(error.data.message || "Error al iniciar sesión");
    });
};
</script>
