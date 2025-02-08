package tcphandler

func (s *RequestHandler) handleTag(request CliRequest) error {
	
	switch request.Command {
	case "view":
		return s.handleTagView(request)
	case "delete":
		return s.handleTagDelete(request)
	default:
		return nil
	}
}

func (s *RequestHandler) handleTagView(request CliRequest) error {
	_, err := s.repo.GetTags()
	if err != nil {
		return err
	}
	// response.Respond(s.Flags.Format, tags)

	return err
}

func (s *RequestHandler) handleTagDelete(request CliRequest) error {
	err := s.repo.DeleteTag(request.Flags.Name)

	return err
}